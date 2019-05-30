# go-templating

Simple text templating package for go programs.

## Functions

Templates are used to make the files somewhat dynamic. Reading and rendering the file is done when GenerateTemplate is called.

Below is a list of the functions available to you and an example of the syntax used.

`hint: it's just golang template syntax`

Function | Description | Example usage
---|---|---
env | Sets the value to an available Environment Variable | {{ env "ENVIRONMENT" }}
default | Use this default value if function fails | {{ default .NonExisting "default value" }}
required | This value must be satisfied or you will get an error | {{ required (default (env "ALWAYS_THERE") "DEFAULT") }}
