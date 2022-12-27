package ui

import "strings"

type Element interface {
	Render() string
	AddClass(s string) Element
	SetId(s string) Element
	GetId() string
	SetAttribute(name string, value string) Element
}

type common struct {
	classes     []string
	attributes  map[string]string
	subElements []*Element
	content     string
}

func (c *common) addElement(s Element) {
	children := &c.subElements
	if children == nil {
		*children = make([]*Element, 0)
	}
	*children = append(*children, &s)
}

func (c *common) addClass(s string) {
	classes := &c.classes
	if classes == nil {
		*classes = make([]string, 0)
	}
	*classes = append(*classes, s)
}

func (c *common) getAttributes() *map[string]string {
	attributes := &c.attributes
	if *attributes == nil {
		*attributes = make(map[string]string)
	}
	return attributes
}

func (c *common) getAttribute(key string) string {
	attributes := c.getAttributes()
	return (*attributes)[key]
}

func (c *common) setAttribute(key string, val string) {
	attributes := c.getAttributes()
	(*attributes)[key] = val
}

func (c *common) getId() string {
	return c.getAttribute("id")
}

func (c *common) setId(s string) {
	c.setAttribute("id", s)
}

func escape(s string) string {
	result := strings.ReplaceAll(s, "\\", "\\\\")
	result = strings.ReplaceAll(result, "\"", "\\\"")
	return result
}

func (c *common) render(tag string) string {
	result := ""
	result += "<" + tag
	if c.classes != nil && len(c.classes) > 0 {
		result += " class=\"" + escape(strings.Join(c.classes, " ")) + "\" "
	}
	if c.attributes != nil && len(c.attributes) > 0 {
		for key, val := range c.attributes {
			if val == "" {
				continue
			}
			result += " " + key + "=\"" + escape(val) + "\" "
		}
	}
	result += ">"
	result += c.content
	if c.subElements != nil {
		for _, s := range c.subElements {
			result += (*s).Render()
		}
	}
	result += "</" + tag + ">"
	return result
}
