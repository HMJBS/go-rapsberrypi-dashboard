// Package fb は Linux FrameBuffer を mmap して RGBA を書き込みます。
package fb

import (
	"fmt"
	"image"
	"os"
	"syscall"
	"unsafe"
)

const (
	fbIOGetVScreenInfo = 0x4600
	fbIOGetFScreenInfo = 0x4602
)

type bitfield struct {
	Offset   uint32
	Length   uint32
	MsbRight uint32
}

type varScreeninfo struct {
	Xres         uint32
	Yres         uint32
	XresVirtual  uint32
	YresVirtual  uint32
	Xoffset      uint32
	Yoffset      uint32
	BitsPerPixel uint32
	Grayscale    uint32
	Red          bitfield
	Green        bitfield
	Blue         bitfield
	Transp       bitfield
	Nonstd       uint32
	Activate     uint32
	Height       uint32
	Width        uint32
	AccelFlags   uint32
	Pixclock     uint32
	LeftMargin   uint32
	RightMargin  uint32
	UpperMargin  uint32
	LowerMargin  uint32
	HsyncLen     uint32
	VsyncLen     uint32
	Sync         uint32
	Vmode        uint32
	Rotate       uint32
	Colorspace   uint32
	Reserved     [4]uint32
}

type fixScreeninfo struct {
	ID           [16]byte
	SmemStart    uintptr
	SmemLen      uint32
	Type         uint32
	TypeAux      uint32
	Visual       uint32
	XpanStep     uint16
	YpanStep     uint16
	YwrapStep    uint16
	LineLength   uint32
	MmioStart    uintptr
	MmioLen      uint32
	Accel        uint32
	Capabilities uint16
	Reserved     [2]uint16
}

// Framebuffer provides a minimal mmap-backed framebuffer writer.
//
// Supported formats:
// - 24bpp packed RGB/BGR (byte-aligned color offsets)
// - 32bpp XRGB8888-like (byte-aligned color offsets)
//
// Other formats return an error.
type Framebuffer struct {
	file *os.File
	mem  []byte

	width      int
	height     int
	lineLength int
	bpp        int

	redIndex   int
	greenIndex int
	blueIndex  int
	alphaIndex int
	bytesPP    int
}

// Open は指定した FrameBuffer デバイスを開きます。
func Open(path string) (*Framebuffer, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	closeOnErr := true
	defer func() {
		if closeOnErr {
			_ = f.Close()
		}
	}()

	var v varScreeninfo
	if err := ioctl(f.Fd(), fbIOGetVScreenInfo, unsafe.Pointer(&v)); err != nil {
		return nil, fmt.Errorf("ioctl FBIOGET_VSCREENINFO: %w", err)
	}
	var fx fixScreeninfo
	if err := ioctl(f.Fd(), fbIOGetFScreenInfo, unsafe.Pointer(&fx)); err != nil {
		return nil, fmt.Errorf("ioctl FBIOGET_FSCREENINFO: %w", err)
	}

	bpp := int(v.BitsPerPixel)
	switch bpp {
	case 24:
	case 32:
	default:
		return nil, fmt.Errorf("unsupported framebuffer bpp=%d", bpp)
	}

	bytesPP := bpp / 8
	lineLength := int(fx.LineLength)
	width := int(v.Xres)
	height := int(v.Yres)

	// Require byte-aligned 8-bit channels.
	if v.Red.Length != 8 || v.Green.Length != 8 || v.Blue.Length != 8 {
		return nil, fmt.Errorf("unsupported color bitfield lengths: r=%d g=%d b=%d", v.Red.Length, v.Green.Length, v.Blue.Length)
	}
	if v.Red.Offset%8 != 0 || v.Green.Offset%8 != 0 || v.Blue.Offset%8 != 0 {
		return nil, fmt.Errorf("unsupported non-byte-aligned color offsets")
	}
	redIndex := int(v.Red.Offset / 8)
	greenIndex := int(v.Green.Offset / 8)
	blueIndex := int(v.Blue.Offset / 8)

	alphaIndex := -1
	if v.Transp.Length == 8 {
		if v.Transp.Offset%8 != 0 {
			return nil, fmt.Errorf("unsupported non-byte-aligned alpha offset")
		}
		alphaIndex = int(v.Transp.Offset / 8)
	}

	// Basic sanity checks.
	maxIndex := redIndex
	if greenIndex > maxIndex {
		maxIndex = greenIndex
	}
	if blueIndex > maxIndex {
		maxIndex = blueIndex
	}
	if alphaIndex > maxIndex {
		maxIndex = alphaIndex
	}
	if maxIndex >= bytesPP {
		return nil, fmt.Errorf("color offsets out of range for %dbpp", bpp)
	}
	if lineLength < width*bytesPP {
		return nil, fmt.Errorf("unexpected line_length=%d for width=%d bytesPP=%d", lineLength, width, bytesPP)
	}

	mem, err := syscall.Mmap(int(f.Fd()), 0, int(fx.SmemLen), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("mmap framebuffer: %w", err)
	}

	closeOnErr = false
	return &Framebuffer{
		file:       f,
		mem:        mem,
		width:      width,
		height:     height,
		lineLength: lineLength,
		bpp:        bpp,
		redIndex:   redIndex,
		greenIndex: greenIndex,
		blueIndex:  blueIndex,
		alphaIndex: alphaIndex,
		bytesPP:    bytesPP,
	}, nil
}

// Close は Framebuffer をクローズし、mmap を解除します。
func (fb *Framebuffer) Close() error {
	if fb == nil {
		return nil
	}
	if fb.mem != nil {
		_ = syscall.Munmap(fb.mem)
		fb.mem = nil
	}
	if fb.file != nil {
		err := fb.file.Close()
		fb.file = nil
		return err
	}
	return nil
}

// Size はフレームバッファの解像度を返します。
func (fb *Framebuffer) Size() image.Point {
	return image.Pt(fb.width, fb.height)
}

// BlitRGBA は RGBA 画像をフレームバッファへ転送します。
func (fb *Framebuffer) BlitRGBA(src *image.RGBA) error {
	if src.Bounds().Dx() != fb.width || src.Bounds().Dy() != fb.height {
		return fmt.Errorf("source size mismatch: got=%dx%d want=%dx%d", src.Bounds().Dx(), src.Bounds().Dy(), fb.width, fb.height)
	}

	// Copy line by line to respect line_length padding.
	for y := 0; y < fb.height; y++ {
		dstLine := fb.mem[y*fb.lineLength : y*fb.lineLength+fb.width*fb.bytesPP]
		srcLine := src.Pix[y*src.Stride : y*src.Stride+fb.width*4]

		di := 0
		for x := 0; x < fb.width; x++ {
			si := x * 4
			r := srcLine[si+0]
			g := srcLine[si+1]
			b := srcLine[si+2]
			a := srcLine[si+3]

			p := dstLine[di : di+fb.bytesPP]
			p[fb.redIndex] = r
			p[fb.greenIndex] = g
			p[fb.blueIndex] = b
			if fb.alphaIndex >= 0 {
				p[fb.alphaIndex] = a
			}

			di += fb.bytesPP
		}
	}
	return nil
}

func ioctl(fd uintptr, req uintptr, arg unsafe.Pointer) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, req, uintptr(arg))
	if errno != 0 {
		return errno
	}
	return nil
}
