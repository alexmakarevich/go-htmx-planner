// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package templs

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func CreateEvent() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h2>Create Event</h2><div><form hx-post=\"/htmx/createEvent\" hx-target=\"#notification\"><input name=\"title\" placeholder=\"title\" required onkeyup=\"this.setCustomValidity(&#39;&#39;) // reset the validation on keyup\" hx-on:htmx:validation:validate=\"if(this.value.length &lt; 1) {\n                    this.setCustomValidity(&#39;The title must be at least 1 char long&#39;) // set the validation error\n                    htmx.find(&#39;#example-form&#39;).reportValidity()          // report the issue\n                }\"> <input type=\"datetime-local\" name=\"date-time\"> <button type=\"submit\">create</button></form></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
