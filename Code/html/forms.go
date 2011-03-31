package main

type Form struct {
	Name, Action, Method string
	Widgets *ComplexPageElement
}

func NewForm(name, action, method string) Form {
	var form Form
	form.Name = name
	form.Action = action
	form.Method = method
	form.Widgets = NewComplexPageElement()
	return form
}

//Adds the given PageElement to this Form with a <br> following it.
func (me *Form) Addln(item PageElement) {
	me.Widgets.Add(item)
	me.Widgets.Add(SimplePageElement("<br>"))
}

//Adds the given PageElement to this Form.
func (me *Form) Add(item PageElement) {
	me.Widgets.Add(item)
}

//Adds a text field with the given label and name to this forum.
func (me *Form) AddTextField(label, name string) {
	me.Widgets.Add(TextField(label, name))
}

//Adds a text field with the given label and name to this forum.
func (me *Form) AddTextFieldln(label, name string) {
	me.Widgets.Add(TextField(label, name) + SimplePageElement("<br>"))
}

//Adds a text field with the given label and name to this forum.
func (me *Form) AddPasswordField(label, name string) {
	me.Widgets.Add(PasswordField(label, name))
}

//Adds a text field with the given label and name to this forum.
func (me *Form) AddPasswordFieldln(label, name string) {
	me.Widgets.Add(PasswordField(label, name) + SimplePageElement("<br>"))
}

//Returns the value of this form as a string.
func (me Form) String() string {
	return "<form " + Attribute("name", me.Name) + Attribute("action", me.Action) + Attribute("method", me.Method) + ">" + me.Widgets.String() + "</from>"
}

//Returns an HTML text field.
func TextField(label, name string) SimplePageElement {
	return SimplePageElement(label + ": <input" + Attribute("type", "text") + Attribute("name", name) + "/>")
}

//Returns an HTML text field.
func PasswordField(label, name string) SimplePageElement {
	return SimplePageElement(label + ": <input" + Attribute("type", "password") + Attribute("name", name) + "/>")
}
