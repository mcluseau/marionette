package marionette

import (
	"encoding/json"
	"fmt"
	"image"
)

type Point struct {
	X float32
	Y float32
}

type Size struct {
	Width  float64
	Height float64
}

type WindowRect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type ElementRect struct {
	Point
	Size
}

type WebElement struct {
	id string //`json:"element-6066-11e4-a52e-4f735466cecf"`
	c  *Client
}

func (e *WebElement) Id() string {
	return e.id
}

func (e *WebElement) GetActiveElement() (*WebElement, error) {
	return e.c.GetActiveElement()
}

func (e *WebElement) FindElement(by By, value string) (*WebElement, error) {
	return e.c.findElement(by, value, &e.id)
}

func (e *WebElement) FindElements(by By, value string) ([]*WebElement, error) {
	return e.c.findElements(by, value, &e.id)
}

func (e *WebElement) Enabled() bool {
	return isElementEnabled(e.c, e.id)
}

func (e *WebElement) Selected() bool {
	return isElementSelected(e.c, e.id)
}

func (e *WebElement) Displayed() bool {
	return isElementDisplayed(e.c, e.id)
}

func (e *WebElement) TagName() string {
	return getElementTagName(e.c, e.id)
}

func (e *WebElement) Text() string {
	return getElementText(e.c, e.id)
}

func Attribute[T any](e *WebElement, name string) (T, error) {
	var out struct {
		Value T `json:"value"`
	}
	err := e.getAttribute(name, &out)
	return out.Value, err
}

func (e *WebElement) getAttribute(name string, dest any) error {
	r, err := e.c.tr.Send("WebDriver:GetElementAttribute", map[string]any{
		"id": e.id, "name": name,
	})
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(r.Value), dest)
}

func (e *WebElement) Attribute(name string) (string, error) {
	return Attribute[string](e, name)
}

func Property[T any](e *WebElement, name string) (T, error) {
	var out struct {
		Value T `json:"value"`
	}
	err := e.getProperty(name, &out)
	return out.Value, err
}

func (e *WebElement) getProperty(name string, dest any) error {
	r, err := e.c.tr.Send("WebDriver:GetElementProperty", map[string]any{
		"id": e.id, "name": name,
	})
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(r.Value), dest)
}

func (e *WebElement) Property(name string) (any, error) {
	return Property[any](e, name)
}

func (e *WebElement) PropertyRaw(name string) (json.RawMessage, error) {
	return Property[json.RawMessage](e, name)
}

func (e *WebElement) PropertyInt(name string) (int, error) {
	return Property[int](e, name)
}

func (e *WebElement) PropertyFloat(name string) (float64, error) {
	return Property[float64](e, name)
}

func (e *WebElement) PropertyString(name string) (string, error) {
	return Property[string](e, name)
}

func (e *WebElement) cssValue(property string, dest any) error {
	r, err := e.c.tr.Send("WebDriver:GetElementCSSValue", map[string]any{
		"id": e.id, "propertyName": property,
	})
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(r.Value), dest)
}

func (e *WebElement) CssValue(property string) (any, error) {
	var out struct {
		Value any `json:"value"`
	}
	err := e.cssValue(property, &out)
	return out.Value, err
}

func (e *WebElement) Rect() (*ElementRect, error) {
	r, err := e.c.tr.Send("WebDriver:GetElementRect", map[string]any{
		"id": e.id,
	})
	if err != nil {
		return nil, err
	}
	d := &ElementRect{}
	err = json.Unmarshal([]byte(r.Value), d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (e *WebElement) Click() {
	r, err := e.c.tr.Send("WebDriver:ElementClick", map[string]any{"id": e.id})
	if err != nil {
		return
	}

	var d = map[string]any{}
	json.Unmarshal([]byte(r.Value), &d)
}

func (e *WebElement) SendKeys(keys string) error {
	//slice := make([]string, 0)
	//for _, v := range keys {
	//	slice = append(slice, fmt.Sprintf("%c", v))
	//}
	//
	//r, err := c.transport.Send("sendKeysToElement", map[string]any{"id": id, "value": slice})
	r, err := e.c.tr.Send("WebDriver:ElementSendKeys", map[string]any{"id": e.id, "text": keys})
	if err != nil {
		return err
	}

	var d = map[string]any{}
	json.Unmarshal([]byte(r.Value), &d)

	return nil
}

func (e *WebElement) Clear() {
	r, err := e.c.tr.Send("WebDriver:ElementClear", map[string]any{"id": e.id})
	if err != nil {
		return
	}

	var d = map[string]any{}
	json.Unmarshal([]byte(r.Value), &d)
}

func (e *WebElement) Location() (*Point, error) {
	r, err := e.Rect()
	if err != nil {
		return nil, err
	}
	return &r.Point, nil
}

func (e *WebElement) Size() (*Size, error) {
	r, err := e.Rect()
	if err != nil {
		return nil, err
	}
	return &r.Size, nil
}

func (e *WebElement) Screenshot() ([]byte, error) {
	id := e.Id()
	return e.c.takeScreenshot(&id)
}

func (e *WebElement) ScreenshotImage() (image.Image, error) {
	id := e.Id()
	return e.c.takeScreenshotImage(&id)
}

func (e *WebElement) UnmarshalJSON(data []byte) error {
	var d map[string]map[string]string
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}
	newId, ok := d["value"][WEBDRIVER_ELEMENT_KEY]
	if !ok {
		return &DriverError{
			ErrorType:  "WebDriverElementKey",
			Message:    fmt.Sprintf("key %v expected in response but not found", WEBDRIVER_ELEMENT_KEY),
			Stacktrace: nil,
		}
	}
	e.id = newId
	return nil
}
