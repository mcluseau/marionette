package marionette

import (
	"context"
	"image"
)

// DownloadImage opens an image URL in the current tab and downloads it.
func (c *Client) DownloadImage(ctx context.Context, src string) (image.Image, error) {
	// TODO: This is a very hack-y way of doing it, but it works more or less.
	//       The other way is to draw image to canvas and transfer data URL,
	//       but some browsers disable this API for security.

	// TODO: Make a new window/tab instead, restore once done.
	_, err := c.Navigate(src)
	if err != nil {
		return nil, err
	}
	el, err := c.FindElement(CSS_SELECTOR, "img")
	if err != nil {
		return nil, err
	}
	width, err := Property[int](el, "naturalWidth")
	if err != nil {
		return nil, err
	}
	height, err := Property[int](el, "naturalHeight")
	if err != nil {
		return nil, err
	}
	for {
		curWidth, err := Property[int](el, "width")
		if err != nil {
			return nil, err
		}
		curHeight, err := Property[int](el, "height")
		if err != nil {
			return nil, err
		}
		if width == curWidth && height == curHeight {
			break
		}
		rect, err := c.GetWindowRect()
		if err != nil {
			return nil, err
		}
		rect.Width += float64(width - curWidth)
		rect.Height += float64(height - curHeight)
		err = c.SetWindowRect(*rect)
		if err != nil {
			return nil, err
		}
	}
	return el.ScreenshotImage()
}
