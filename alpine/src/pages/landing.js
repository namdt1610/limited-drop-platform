import { LandingLogic } from "./landing.logic.js";

customElements.define(
  "landing-page",
  class extends HTMLElement {
    connectedCallback() {
      // Render khung xương ngay lập tức, component sống, khác với constructor chỉ khai báo trên RAM
      this.innerHTML = /*html*/ `
            <div x-data="landingTeaser" 
                 @mousemove.window="updateCursor($event)"
                 class="landing-hover min-h-screen text-foreground cursor-none">
                
                <!-- Custom Cursor -->
                <div class="fixed w-8 h-8 border border-primary rounded-full pointer-events-none z-[9999] transition-transform duration-100 ease-out mix-blend-difference hidden md:block"
                     :style="'left: ' + cursorX + 'px; top: ' + cursorY + 'px; transform: translate(-50%, -50%) scale(' + (isHovering ? 2 : 1) + ')'">
                    <div class="absolute inset-0 bg-primary/20 rounded-full animate-pulse"></div>
                </div>

                <div class="relative z-10 min-h-screen flex flex-col">
                    <teaser-hero-section></teaser-hero-section>
                    <teaser-story-section></teaser-story-section>
                    <system-stories-section></system-stories-section>
                    <cinematic-video-section></cinematic-video-section>
                    <teaser-footer></teaser-footer>
                </div>

                <!-- CTA: centered fixed at bottom -->
                <div class="fixed bottom-3 sm:bottom-6 left-1/2 -translate-x-1/2 z-30 flex justify-center w-full pointer-events-none px-4">
                    <div class="pointer-events-auto group"
                         @mousemove="handleMagnetic($event)"
                         @mouseleave="resetMagnetic"
                         :style="magneticStyle">
                        <a href="#"
                           @click.prevent="if (canEnter()) $dispatch('route', { page: 'drop', id: activeDropId })"
                           :aria-disabled="!canEnter()"
                           :class="canEnter() ? 'inline-flex items-center gap-2 px-6 sm:px-12 py-3 sm:py-6 text-xs sm:text-sm font-semibold shadow-lg transition chrome-btn chrome-btn--breathing' : 'inline-flex items-center gap-2 px-6 sm:px-12 py-3 sm:py-6 text-xs sm:text-sm font-semibold shadow-lg transition chrome-btn opacity-60 cursor-not-allowed'"
                           x-text="canEnter() ? 'Vào Đấu Trường' : 'Đấu Trường sắp mở'"></a>
                    </div>
                </div>
            </div>
        `;

      // Đăng ký logic tách riêng (guarded để chỉ register 1 lần)
      if (!Alpine.store("landingTeaserInitialized")) {
        Alpine.data("landingTeaser", LandingLogic);
        Alpine.store("landingTeaserInitialized", true);
      }
    }
  }
);
