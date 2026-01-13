// Mercury Background - Galaxy Starfield Web Component
class MercuryBackground extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    this.time = 0;
    this.animationFrameId = null;
    this.stars = [];
    this.animate = this.animate.bind(this);
  }

  connectedCallback() {
    this.generateStars();
    this.render();
    this.animationFrameId = requestAnimationFrame(this.animate);
  }

  disconnectedCallback() {
    if (this.animationFrameId) {
      cancelAnimationFrame(this.animationFrameId);
    }
  }

  generateStars() {
    this.stars = [];
    for (let i = 0; i < 200; i++) {
      this.stars.push({
        x: Math.random() * 100,
        y: Math.random() * 100,
        size: Math.random() * 2.5,
        opacity: Math.random() * 0.7 + 0.3,
        twinkleSpeed: Math.random() * 0.02 + 0.008,
        twinklePhase: Math.random() * Math.PI * 2,
        moveSpeedX: (Math.random() - 0.5) * 0.01,
        moveSpeedY: (Math.random() - 0.5) * 0.015,
      });
    }
  }

  animate() {
    this.time += 0.016;
    this.updateStars();
    this.animationFrameId = requestAnimationFrame(this.animate);
  }

  updateStars() {
    const starsContainer = this.shadowRoot.querySelector(".stars");
    if (!starsContainer) return;

    this.stars.forEach((star, index) => {
      const starEl = starsContainer.children[index];
      if (starEl) {
        // Update position with movement
        star.x = (star.x + star.moveSpeedX + 100) % 100;
        star.y = (star.y + star.moveSpeedY + 100) % 100;

        // Update twinkle
        const twinkle =
          Math.sin(this.time * star.twinkleSpeed + star.twinklePhase) * 0.4 +
          0.6;
        starEl.style.opacity = (star.opacity * twinkle).toString();

        // Update position
        starEl.style.left = `${star.x}%`;
        starEl.style.top = `${star.y}%`;
      }
    });
  }

  render() {
    const starsHTML = this.stars
      .map(
        (star) =>
          `<div class="star" style="left: ${star.x}%; top: ${star.y}%; width: ${star.size}px; height: ${star.size}px; opacity: ${star.opacity};"></div>`
      )
      .join("");

    this.shadowRoot.innerHTML = `
      <style>
        :host {
          position: fixed;
          inset: 0;
          z-index: -1;
          pointer-events: none;
        }

        .galaxy-bg {
          position: absolute;
          inset: 0;
          background: radial-gradient(
            ellipse at 40% 30%,
            hsl(220 30% 15%) 0%,
            hsl(220 16% 8%) 40%,
            hsl(220 15% 5%) 100%
          );
        }

        .stars-layer {
          position: absolute;
          inset: 0;
        }

        .stars {
          position: absolute;
          inset: 0;
          width: 100%;
          height: 100%;
        }

        .star {
          position: absolute;
          background: radial-gradient(circle, rgba(255, 255, 255, 0.9) 0%, rgba(255, 255, 255, 0.3) 70%, transparent 100%);
          border-radius: 50%;
          box-shadow: 0 0 2px rgba(100, 150, 255, 0.6), 0 0 4px rgba(100, 150, 255, 0.3);
          will-change: opacity;
        }

        /* Distant nebula glow */
        .nebula {
          position: absolute;
          border-radius: 50%;
          filter: blur(60px);
          opacity: 0.1;
        }

        .nebula-1 {
          top: 20%;
          left: 15%;
          width: 300px;
          height: 300px;
          background: radial-gradient(circle, rgba(100, 150, 255, 0.3) 0%, transparent 70%);
          animation: float 30s ease-in-out infinite;
        }

        .nebula-2 {
          bottom: 10%;
          right: 10%;
          width: 250px;
          height: 250px;
          background: radial-gradient(circle, rgba(150, 100, 255, 0.2) 0%, transparent 70%);
          animation: float 35s ease-in-out infinite 5s;
        }

        @keyframes float {
          0%, 100% {
            transform: translate(0, 0);
          }
          25% {
            transform: translate(20px, -20px);
          }
          50% {
            transform: translate(-10px, 10px);
          }
          75% {
            transform: translate(15px, 15px);
          }
        }
      </style>

      <div class="galaxy-bg"></div>
      <div class="stars-layer">
        <div class="nebula nebula-1"></div>
        <div class="nebula nebula-2"></div>
        <div class="stars">${starsHTML}</div>
      </div>
    `;
  }
}

// Register the custom element
customElements.define("mercury-background", MercuryBackground);
