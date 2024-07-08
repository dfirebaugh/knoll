class MyCard extends HTMLElement {
  constructor() {
    super();
    const shadow = this.attachShadow({ mode: "open" });

    const container = document.createElement("div");
    container.setAttribute("class", "card");

    const title = document.createElement("h2");
    title.textContent = this.getAttribute("title") || "Card Title";

    const content = document.createElement("p");
    content.textContent = this.getAttribute("content") || "Card Content";

    const style = document.createElement("style");
    style.textContent = `
            .card {
                background-color: lightgray;
                border: 1px solid #ccc;
                padding: 16px;
                border-radius: 8px;
                box-shadow: 0 2px 5px rgba(0,0,0,0.1);
            }
            h2 {
                margin: 0;
                font-size: 1.5em;
                color: #333;
            }
            p {
                font-size: 1em;
                color: #666;
            }
        `;

    shadow.appendChild(style);
    shadow.appendChild(container);
    container.appendChild(title);
    container.appendChild(content);
  }

  static get observedAttributes() {
    return ["title", "content"];
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (name === "title") {
      this.shadowRoot.querySelector("h2").textContent = newValue;
    } else if (name === "content") {
      this.shadowRoot.querySelector("p").textContent = newValue;
    }
  }
}

customElements.define("my-card", MyCard);
