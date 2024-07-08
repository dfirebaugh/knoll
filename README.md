# Knoll

[demo](https://dfirebaugh.github.io/knoll/?element=purple+button)

Knoll is a minimalistic gallery for custom html elements which can be used to
help you develop out of band from your main application or as a tool for
documentation.

### Installation

```sh
go install github.com/dfirebaugh/knoll@latest
```

> make sure `~/go/bin` is in your path

### Usage

1. Create a YAML configuration file (e.g., `config.yaml`) that defines your
   custom elements. Example config below.

2. Run Knoll to generate the static site:
   ```sh
   knoll --output output config.yaml
   ```

3. Optionally, serve the generated static site:
   ```sh
   knoll --serve --port 8080 --output output config.yaml
   ```

   This will start a web server at `http://localhost:8080` serving the generated
   files from the `output` directory.

### Example config

```yaml
   scripts:
     - examples/button.js

   elements:
     - name: "purple button"
       tag: "my-button"
       attributes:
         - name: "label"
           type: "string"
           default: "Click Me"
         - name: "color"
           type: "string"
           default: "white"
         - name: "background-color"
           type: "string"
           default: "blue"
       exampleData:
         label: "Submit"
         color: "white"
         background-color: "purple"
       script: "examples/button.js"
       properties:
         - name: "disabled"
           type: "boolean"
           default: false
       events:
         - name: "button-click"
           description: "Triggered when the button is clicked."

     - name: "green button"
       tag: "my-button"
       attributes:
         - name: "label"
           type: "string"
           default: "Click Me"
         - name: "color"
           type: "string"
           default: "white"
         - name: "background-color"
           type: "string"
           default: "blue"
       exampleData:
         label: "Test Button"
         color: "white"
         background-color: "green"
       script: "examples/button.js"
       properties:
         - name: "disabled"
           type: "boolean"
           default: false
       events:
         - name: "button-click"
           description: "Triggered when the button is clicked."

     - name: "card"
       tag: "my-card"
       attributes:
         - name: "title"
           type: "string"
           default: "Card Title"
         - name: "content"
           type: "string"
           default: "Card Content"
       exampleData:
         title: "Welcome"
         content: "This is an example card"
       script: "examples/card.js"
       properties: []
       events: []
```
