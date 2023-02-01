package processing

import (
	"context"

	"github.com/imgproxy/imgproxy/v3/config"
	"github.com/imgproxy/imgproxy/v3/imagedata"
	"github.com/imgproxy/imgproxy/v3/imath"
	"github.com/imgproxy/imgproxy/v3/options"
	"github.com/imgproxy/imgproxy/v3/vips"
)

var watermarkPipeline = pipeline{
	prepare,
	scaleOnLoad,
	importColorProfile,
	scale,
	rotateAndFlip,
	padding,
}

func prepareWatermark(wm *vips.Image, wmData *imagedata.ImageData, opts *options.WatermarkOptions, imgWidth, imgHeight, framesCount int) error {
	if err := wm.Load(wmData, 1, 1.0, 1); err != nil {
		return err
	}

	po := options.NewProcessingOptions()
	po.ResizingType = options.ResizeFit
	po.Dpr = 1
	po.Enlarge = true
	po.Format = wmData.Type

	if opts.Scale > 0 {
		po.Width = imath.Max(imath.Scale(imgWidth, opts.Scale), 1)
		po.Height = imath.Max(imath.Scale(imgHeight, opts.Scale), 1)
	}

	if opts.Replicate {
		po.Padding.Enabled = true
		po.Padding.Left = int(opts.Gravity.X / 2)
		po.Padding.Right = int(opts.Gravity.X) - po.Padding.Left
		po.Padding.Top = int(opts.Gravity.Y / 2)
		po.Padding.Bottom = int(opts.Gravity.Y) - po.Padding.Top
	}

	if err := watermarkPipeline.Run(context.Background(), wm, po, wmData); err != nil {
		return err
	}

	if opts.Replicate || framesCount > 1 {
		// We need to copy image if we're going to replicate.
		// Replication requires image to be read several times, and this requires
		// random access to pixels
		if err := wm.CopyMemory(); err != nil {
			return err
		}
	}

	if opts.Replicate {
		if err := wm.Replicate(imgWidth, imgHeight); err != nil {
			return err
		}
	} else {
		left, top := calcPosition(imgWidth, imgHeight, wm.Width(), wm.Height(), &opts.Gravity, true)
		if err := wm.Embed(imgWidth, imgHeight, left, top); err != nil {
			return err
		}
	}

	if framesCount > 1 {
		if err := wm.Replicate(imgWidth, imgWidth*framesCount); err != nil {
			return err
		}
	}

	return nil
}

func applyWatermark(img *vips.Image, wmData *imagedata.ImageData, opts *options.WatermarkOptions, framesCount int) error {
	if err := img.RgbColourspace(); err != nil {
		return err
	}

	wm := new(vips.Image)
	defer wm.Clear()

	width := img.Width()
	height := img.Height()

	if err := prepareWatermark(wm, wmData, opts, width, height/framesCount, framesCount); err != nil {
		return err
	}

	opacity := opts.Opacity * config.WatermarkOpacity

	return img.ApplyWatermark(wm, opacity)
}

func watermark(pctx *pipelineContext, img *vips.Image, po *options.ProcessingOptions, imgdata *imagedata.ImageData) error {
	if !po.Watermark.Enabled || imagedata.Watermark == nil {
		return nil
	}

	return applyWatermark(img, imagedata.Watermark, &po.Watermark, 1)
}
