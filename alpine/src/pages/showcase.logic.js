export function ShowcaseLogic() {
  return {
    products: [
      {
        id: 101,
        name: "CHROME OVERDRIVE JACKET",
        price: 4500000,
        thumbnail: "https://images.unsplash.com/photo-1551028719-00167b16eac5?auto=format&fit=crop&w=800&q=80",
        category: "Outerwear",
        slug: "chrome-overdrive-jacket"
      },
      {
        id: 102,
        name: "NEON PULSE SNEAKERS",
        price: 3200000,
        thumbnail: "https://images.unsplash.com/photo-1542291026-7eec264c27ff?auto=format&fit=crop&w=800&q=80",
        category: "Footwear",
        slug: "neon-pulse-sneakers"
      },
      {
        id: 103,
        name: "VOID CARGO PANTS",
        price: 2800000,
        thumbnail: "https://images.unsplash.com/photo-1552902865-b72c031ac5ea?auto=format&fit=crop&w=800&q=80",
        category: "Bottoms",
        slug: "void-cargo-pants"
      },
      {
        id: 104,
        name: "CYBER GLASSES V1",
        price: 1500000,
        thumbnail: "https://images.unsplash.com/photo-1511499767390-90342f16b117?auto=format&fit=crop&w=800&q=80",
        category: "Accessories",
        slug: "cyber-glasses-v1"
      },
      {
        id: 105,
        name: "TECH-WEAR BACKPACK",
        price: 3800000,
        thumbnail: "https://images.unsplash.com/photo-1553062407-98eeb64c6a62?auto=format&fit=crop&w=800&q=80",
        category: "Accessories",
        slug: "tech-wear-backpack"
      },
      {
        id: 106,
        name: "SYNTH WAVE TEE",
        price: 950000,
        thumbnail: "https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?auto=format&fit=crop&w=800&q=80",
        category: "Tops",
        slug: "synth-wave-tee"
      },
      {
        id: 107,
        name: "GLITCH BEANIE",
        price: 650000,
        thumbnail: "https://images.unsplash.com/photo-1576871337622-98d48d365da2?auto=format&fit=crop&w=800&q=80",
        category: "Accessories",
        slug: "glitch-beanie"
      },
      {
        id: 108,
        name: "CARBON FIBER WATCH",
        price: 12000000,
        thumbnail: "https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&w=800&q=80",
        category: "Accessories",
        slug: "carbon-fiber-watch"
      }
    ],
    loading: false,
    error: null,

    init() {
      this.initRevealObserver();
    },

    initRevealObserver() {
      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              entry.target.classList.add("reveal-visible");
            }
          });
        },
        { threshold: 0.1 }
      );
      setTimeout(() => {
        document
          .querySelectorAll(".reveal")
          .forEach((el) => observer.observe(el));
      }, 500);
    },

    openProduct(p) {
      window.dispatchEvent(
        new CustomEvent("route", { 
          detail: { page: "product", id: p.id || p.slug } 
        })
      );
    },
  };
}
