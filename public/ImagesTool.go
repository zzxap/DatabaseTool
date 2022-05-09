package public

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"net/http"
	//"mime/multipart"
	//"image/color"
	"math/rand"
	"os"

	"github.com/disintegration/imaging"
	//"runtime"

	"strconv"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

// 改变大小
func ImageFile_resize(file *os.File, toPath string, width uint, height uint) {
	// 打开图片并解码
	Log("topath=" + toPath)
	if strings.Contains(strings.ToLower(toPath), ".png") {
		ImageFile_resize_png(file, toPath, width, height)

	} else {
		ImageFile_resize_jpg(file, toPath, width, height)
	}

}
func ImageFile_resize_jpg(file *os.File, toPath string, width uint, height uint) {
	// 打开图片并解码
	Log("resize jpg 12")
	if file == nil {
		Log("file is nil")
		return
	}
	img, err := jpeg.Decode(file)
	Log("resize jpg 13")
	if err != nil {
		Log("jpeg.Decode error=")
		Log(err)
		return
	}
	//file.Close()
	Log("resize jpg 14")
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(200, 0, img, resize.Lanczos3)
	Log("resize jpg 15")
	out, err := os.Create(toPath)
	if err != nil {
		Log(err)
	}
	defer out.Close()
	Log("resize jpg 15")
	// write new image to file
	jpeg.Encode(out, m, nil)

}

//压缩图片
func CompressImage(frompath string, topath string, width int) error {

	src, err := imaging.Open(frompath)
	if err != nil {
		Log("failed to open image: %v", err)
		return err
	}

	// Resize the cropped image to width = 200px preserving the aspect ratio.
	src = imaging.Resize(src, width, 0, imaging.Lanczos)

	// Save the resulting image as JPEG.
	err = imaging.Save(src, topath)
	if err != nil {
		Log("failed to save image: %v", err)
		return err
	}
	return nil
	// write new image to file

}

//压缩图片
func CompressTwoImage(frompath string, bigpath string, smallpath string, logopath string, bigwidth int, smallwidth int) error {

	src, err := imaging.Open(frompath)
	if err != nil {
		Log("failed to open image: %v", err)
		return err
	}

	// Resize the cropped image to width = 200px preserving the aspect ratio.
	bigimg := imaging.Resize(src, bigwidth, 0, imaging.Lanczos)
	err = imaging.Save(bigimg, bigpath)
	if err != nil {
		Log("failed to save image: %v", err)
		return err
	}
	smallimg := imaging.Resize(src, smallwidth, 0, imaging.Lanczos)
	err = imaging.Save(smallimg, smallpath)

	if err != nil {
		Log("failed to save image: %v", err)
		return err
	}
	return nil

}

//压缩图片
func CompressImage2(frompath string, topath string, width uint) error {

	PthSep := string(os.PathSeparator)
	array := strings.Split(frompath, PthSep)
	filename := array[len(array)-1]
	//read, err := ioutil.ReadFile(frompath)
	file, err := os.Open(frompath)
	if err != nil {
		Log(err)
		return err
	}

	buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
	if _, err = file.Read(buff); err != nil {
		fmt.Println(err) // do something with that error
		return err
	}
	ContentType := http.DetectContentType(buff)
	Log("ContentType=" + ContentType)
	//fmt.Println("upload picture 10")
	// decode jpeg into image.Image
	var img image.Image
	var errr error
	if strings.Contains(strings.ToLower(ContentType), ".png") {
		fmt.Printf("is png")
		img, errr = png.Decode(file)

	} else {
		fmt.Printf("is jpg")
		//img, errr = jpeg.Decode(file)
		img, _, errr = image.Decode(file)

	}

	if errr != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", filename, errr)
		return errr
	}

	file.Close()

	//生成大图片------------------------------
	bigImg := resize.Resize(width, 0, img, resize.Lanczos3)

	bigout, err := os.Create(topath)

	if err != nil {
		Log(err)
		return err
	}
	defer bigout.Close()
	if strings.Contains(strings.ToLower(ContentType), ".png") {
		return png.Encode(bigout, bigImg)

	} else {
		return jpeg.Encode(bigout, bigImg, &jpeg.Options{100})

	}

	// write new image to file

}

//压缩图片
func CompressTwoImage22(frompath string, bigpath string, smallpath string, logopath string, bigwidth uint, smallwidth uint) {

	PthSep := string(os.PathSeparator)
	array := strings.Split(frompath, PthSep)
	filename := array[len(array)-1]
	//fmt.Printf("frompath=" + frompath)
	file, err := os.Open(frompath)
	if err != nil {
		Log(err)
	}
	//fmt.Println("upload picture 10")
	// decode jpeg into image.Image
	var img image.Image
	var errr error
	if strings.Contains(strings.ToLower(filename), ".png") {
		fmt.Printf("is png")
		img, errr = png.Decode(file)
		if errr != nil {
			//fmt.Printf("is png errorrrr  \n")
			//Log(errr)
			img, errr = jpeg.Decode(file)
			if errr != nil {
				fmt.Printf("is jpg error \n ")
				//Log(errr)

			}
		}

	} else {
		fmt.Printf("is jpg")
		img, errr = jpeg.Decode(file)
		if errr != nil {
			fmt.Printf("is jpg errorr")
			//Log(errr)
			img, errr = png.Decode(file)
		}
	}
	//fmt.Printf("make big pic111  ")
	if errr != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", filename, errr)
		return
	}

	file.Close()
	fmt.Printf("make big pic222 \n")
	//生成大图片------------------------------
	bigImg := resize.Resize(bigwidth, 0, img, resize.Lanczos3)
	fmt.Printf("make big pic223 \n")
	bigout, err := os.Create(bigpath)
	fmt.Printf("make big pic224 \n")
	if err != nil {
		Log(err)
	}
	defer bigout.Close()

	// write new image to file
	jpeg.Encode(bigout, bigImg, &jpeg.Options{jpeg.DefaultQuality})
	fmt.Printf("make small pic")
	//生成小图片------------------------------
	bigImg = resize.Resize(smallwidth, 0, img, resize.Lanczos3)

	bigout, err = os.Create(smallpath)

	if err != nil {
		Log(err)
	}
	defer bigout.Close()

	// write new image to file
	jpeg.Encode(bigout, bigImg, nil)

}

//水印
func SignImage(frompath string, logopath string) {

	topath := strings.Replace(frompath, ".png", "temp.png", -1)
	topath = strings.Replace(frompath, ".jpg", "temp.jpg", -1)
	topath = strings.Replace(frompath, ".jpeg", "temp.jpeg", -1)
	topath = strings.Replace(frompath, ".PNG", "temp.PNG", -1)
	topath = strings.Replace(frompath, ".JPG", "temp.JPG", -1)
	topath = strings.Replace(frompath, ".JPEG", "temp.JPEG", -1)

	imgb, erropen := os.Open(logopath)
	if erropen != nil {
		Log(erropen)
	}
	watermark, errdecode := png.Decode(imgb)
	if errdecode != nil {
		Log(errdecode)
	}
	defer imgb.Close()

	wmb, _ := os.Open(frompath)
	var sourceImage image.Image
	var err error
	if strings.Contains(strings.ToLower(frompath), ".png") {
		//Log("is png")
		sourceImage, err = png.Decode(wmb)
		if err != nil {
			Log(err)
			return
			//sourceImage, err = jpeg.Decode(wmb)
		}
	} else {
		//Log("is jpg")
		sourceImage, err = jpeg.Decode(wmb)
		if err != nil {
			Log(err)
			return
			//sourceImage, err = png.Decode(wmb)
		}
	}
	if err != nil {
		Log(err)
	}
	defer wmb.Close()

	//把水印写到右下角，并向0坐标各偏移10个像素
	offset := image.Pt(sourceImage.Bounds().Dx()-watermark.Bounds().Dx()-10, sourceImage.Bounds().Dy()-watermark.Bounds().Dy()-10)
	b := sourceImage.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, sourceImage, image.ZP, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

	imgw, err := os.Create(topath)
	jpeg.Encode(imgw, m, &jpeg.Options{jpeg.DefaultQuality})
	imgw.Close()
	if err == nil {
		//删除替换原来的文件
		os.Remove(frompath)
		os.Rename(topath, frompath)
	}

}
func ImageFile_resize_png(file *os.File, toPath string, width uint, height uint) {

	Log("resize png 12")
	// decode jpeg into image.Image
	img, err := png.Decode(file)
	if err != nil {
		Log(err)
		return
	}
	//defer file.Close()
	Log("resize png 13")
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(width, 0, img, resize.Lanczos3)
	Log("resize png 14")
	out, err := os.Create(toPath)
	if err != nil {
		Log(err)
		return
	}
	defer out.Close()
	Log("resize png 15")
	// write new image to file
	png.Encode(out, m)

}

// 改变大小
func Image_resize(fromPath string, toPath string, width uint, height uint) {
	// 打开图片并解码
	fromPath = strings.Replace(fromPath, "/", "\\", -1)
	toPath = strings.Replace(toPath, "/", "\\", -1)

	if strings.Contains(strings.ToLower(toPath), ".png") {
		Image_resize_png(fromPath, toPath, width, height)

	} else {
		Image_resize_jpg(fromPath, toPath, width, height)
	}

	//jpeg.Encode(file_out, canvas, &jpeg.Options{80})

	// cmd_watermark(to, strings.Replace(to, ".jpg", "@big.jpg", 1))
	//Image_thumbnail(to, 200, 200)
}

func Image_resize_png(fromPath string, toPath string, width uint, height uint) {
	fmt.Println("upload png 10")
	// 打开图片并解码
	ff, errr := os.Open(fromPath)
	if errr != nil {
		Log(errr)
		return
	}

	img, er := png.Decode(ff)
	if er != nil {
		Log(er)
		return
	}
	ff.Close()
	Log("resize png 13")
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(200, 0, img, resize.Lanczos3)
	Log("resize png 14")
	out, errrr := os.Create(toPath)
	if errrr != nil {
		Log(errrr)
	}
	defer out.Close()
	Log("resize png 15")
	// write new image to file
	png.Encode(out, m)

}

func Image_resize_jpg(from string, to string, width uint, height uint) {
	fmt.Println("upload jpg =" + from)
	file, err := os.Open(from)
	if err != nil {
		Log(err)
	}
	fmt.Println("upload picture 10")
	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		Log(err)
	}
	file.Close()
	fmt.Println("upload picture 11")
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(200, 0, img, resize.Lanczos3)

	out, err := os.Create(to)
	if err != nil {
		Log(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

}

// 水印
func Image_watermark(file string, to string) {
	// 打开图片并解码
	file_origin, _ := os.Open(file)
	origin, _ := jpeg.Decode(file_origin)
	defer file_origin.Close()
	// 打开水印图并解码
	file_watermark, _ := os.Open("watermark.png")
	watermark, _ := png.Decode(file_watermark)
	defer file_watermark.Close()
	//原始图界限
	origin_size := origin.Bounds()
	//创建新图层
	canvas := image.NewNRGBA(origin_size)
	// 贴原始图
	draw.Draw(canvas, origin_size, origin, image.ZP, draw.Src)
	// 贴水印图
	draw.Draw(canvas, watermark.Bounds().Add(image.Pt(30, 30)), watermark, image.ZP, draw.Over)
	//生成新图片
	create_image, _ := os.Create(to)
	jpeg.Encode(create_image, canvas, &jpeg.Options{95})
	defer create_image.Close()
}

// 随机生成文件名
func Random_name() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Int())
}
