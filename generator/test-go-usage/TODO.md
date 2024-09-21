from vscode ext :

```go
package main

struct Data {
    message string
}

gotempl (d Data) Method(greeting string) {
  {{ // comment }}
	// { d.message }
	// { greeting }
}

gotempl test(comp templ.Component) {
	<div>
		@comp {
			<div>Children</div>
		}
	</div>
}

gotempl Hello[T ~string]() {
	@Data{
		message: "You can implement methods on a type.",
	}.Method("hello") {
		@test(Data{message: "You can implement methods on a type."}.Method("hello"))
		<div>
			{ children... }
		</div>
	}
}

```
