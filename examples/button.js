class MyButton extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    this.render();
  }

  static get observedAttributes() {
    return ["label", "color", "background-color"];
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (oldValue !== newValue) {
      this.render();
    }
  }

  get disabled() {
    return this.hasAttribute("disabled");
  }

  set disabled(value) {
    if (value) {
      this.setAttribute("disabled", "");
    } else {
      this.removeAttribute("disabled");
    }
    this.render();
  }

  render() {
    this.shadowRoot.innerHTML = `
      <style>
        button {
          background-color: ${this.getAttribute("background-color") || "blue"};
          color: ${this.getAttribute("color") || "white"};
          padding: 10px 20px;
          border: none;
          border-radius: 5px;
          cursor: pointer;
          font-size: 16px;
        }
        button:hover {
          background-color: #0056b3;
        }
        button:disabled {
          background-color: grey;
          cursor: not-allowed;
        }
      </style>
      <button ${this.disabled ? "disabled" : ""}>${
      this.getAttribute("label") || "Click me"
    }</button>
    `;

    const button = this.shadowRoot.querySelector("button");
    button.addEventListener("click", () => {
      this.dispatchEvent(
        new CustomEvent("button-click", {
          detail: { message: "Button was clicked!" },
          bubbles: true,
          composed: true,
        }),
      );
    });
  }
}

customElements.define("my-button", MyButton);
