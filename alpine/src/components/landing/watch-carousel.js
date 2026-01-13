// Watch Carousel Component for Alpine.js
document.addEventListener("alpine:init", () => {
  Alpine.data("watchCarousel", () => ({
    currentIndex: 0,
    images: [
      "imgs/watch-1.webp",
      "imgs/watch-2.webp",
      "imgs/watch-3.webp",
      "imgs/watch-4.webp",
      "imgs/watch-5.webp",
    ],
    autoPlayInterval: null,
    isHovered: false,
    mousePos: { x: 0, y: 0 },

    init() {
      // Start auto-play
      this.startAutoPlay();
    },

    startAutoPlay() {
      if (this.images.length <= 1) return;
      this.autoPlayInterval = setInterval(() => {
        this.next();
      }, 3500);
    },

    stopAutoPlay() {
      if (this.autoPlayInterval) {
        clearInterval(this.autoPlayInterval);
        this.autoPlayInterval = null;
      }
    },

    onMouseMove(e, el) {
      const rect = el.getBoundingClientRect();
      const x = (e.clientX - rect.left - rect.width / 2) / rect.width;
      const y = (e.clientY - rect.top - rect.height / 2) / rect.height;
      this.mousePos = { x: x * 10, y: y * 10 };
    },

    onMouseLeave() {
      this.isHovered = false;
      this.mousePos = { x: 0, y: 0 };
    },

    next() {
      this.currentIndex = (this.currentIndex + 1) % this.images.length;
    },

    prev() {
      this.currentIndex =
        (this.currentIndex - 1 + this.images.length) % this.images.length;
    },

    goTo(index) {
      this.currentIndex = index;
    },
  }));
});

// Register watch-carousel as a custom element with inline template
customElements.define(
  "watch-carousel",
  class extends HTMLElement {
    connectedCallback() {
      this.innerHTML = /*html*/ `
    <div class="relative aspect-square overflow-hidden rounded-2xl" x-data="watchCarousel" @mousemove="onMouseMove($event, $el)" @mouseenter="isHovered = true" @mouseleave="onMouseLeave()">
        <!-- Glow effect - intensifies on hover -->
        <div class="absolute inset-0 bg-primary/10 rounded-3xl blur-2xl transition-all duration-500" :class="isHovered ? 'opacity-100 scale-125' : 'opacity-70 scale-110'"></div>

        <!-- Main Image Container -->
        <div class="relative w-full h-full rounded-3xl overflow-hidden border bg-card transition-all duration-500" :class="isHovered ? 'border-blue-500/30 shadow-2xl' : 'theme-border shadow-xl'">
            <template x-for="(image, index) in images" :key="index">
                <div :class="index === currentIndex ? 'opacity-100 scale-100 absolute inset-0 transition-all duration-1000 ease-out-expo' : 'opacity-0 scale-105 absolute inset-0 transition-all duration-1000 ease-out-expo'">
                    <img
                        :src="image"
                        :alt="'Watch ' + (index + 1)"
                        class="w-full h-full object-cover transition-transform duration-700"
                        :class="isHovered ? 'scale-105' : ''"
                        loading="lazy"
                    />
                </div>
            </template>

            <!-- Shine overlay on hover -->
            <div :class="isHovered ? 'absolute inset-0 bg-gradient-to-tr from-transparent via-white/10 to-transparent transition-opacity duration-500 pointer-events-none opacity-100' : 'absolute inset-0 bg-gradient-to-tr from-transparent via-white/10 to-transparent transition-opacity duration-500 pointer-events-none opacity-0'"></div>
        </div>

        <!-- Navigation Arrows -->
        <button
            @click="prev()"
            class="carousel-control left-4"
        >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
            </svg>
        </button>

        <button
            @click="next()"
            class="carousel-control right-4"
        >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
            </svg>
        </button>

        <!-- Indicators -->
        <div class="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
            <template x-for="(image, index) in images" :key="index">
                <button
                    @click="goTo(index)"
                    class="carousel-dot"
                    :class="{ 'carousel-dot--active': index === currentIndex, 'carousel-dot--inactive': index !== currentIndex }"
                ></button>
            </template>
        </div>
    </div>
`;
    }
  }
);
