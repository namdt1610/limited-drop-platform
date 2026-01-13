// Define teaser-hero-section directly as a custom element
customElements.define(
  "teaser-hero-section",
  class extends HTMLElement {
    connectedCallback() {
      this.innerHTML = /*html*/ `
        <main x-data="{ heroParallax: 0 }" 
              @scroll.window="heroParallax = window.pageYOffset * 0.1"
              class="flex-1 flex items-center py-12 px-6 sm:px-8 md:px-12 relative overflow-hidden">
            
            <!-- Background Decorative Text -->
            <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-[20vw] font-black text-white/[0.02] select-none pointer-events-none whitespace-nowrap"
                 :style="'transform: translate(-50%, -50%) translateX(' + heroParallax + 'px)'">
                SYMBIOSIS CHROME
            </div>

            <div class="max-w-7xl mx-auto w-full relative z-10">
                <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-16 xl:gap-24 items-center">
                    <!-- Content -->
                    <div class="flex flex-col justify-center space-y-8 text-center lg:text-left order-2 lg:order-1">
                        <!-- Badge - with hover pulse -->
                        <div class="reveal inline-flex items-center gap-2 px-4 py-2 rounded-full bg-card/50 border theme-border w-fit mx-auto lg:mx-0 transition-all duration-300 hover:bg-card/70 hover:scale-105 hover:shadow-lg hover:shadow-accent/15 cursor-default backdrop-blur">
                            <svg class="w-4 h-4 text-foreground animate-pulse" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"></path>
                            </svg>
                            <span class="text-xs font-medium tracking-[0.14em] text-foreground uppercase">
                                Bộ sưu tập Symbiote
                            </span>
                        </div>

                        <!-- Story Arc - The Promise -->
                        <div class="reveal reveal-delay-1 space-y-6 group">
                            <!-- Opening Statement -->
                            <div class="space-y-3">
                                <p class="text-sm sm:text-base text-muted-foreground uppercase tracking-[0.15em] transition-colors group-hover:text-primary">
                                    Giao thoa Công nghệ & Nghệ thuật
                                </p>
                                <h1 class="font-display tracking-tight">
                                    <span class="block text-4xl sm:text-5xl md:text-6xl lg:text-6xl xl:text-7xl font-black uppercase leading-[0.93] tracking-[-0.02em] text-foreground transition-all duration-500 group-hover:tracking-tighter">
                                        SYMBIOSIS CHROME
                                    </span>
                                    <span class="block text-2xl sm:text-3xl lg:text-4xl font-bold text-muted-foreground uppercase tracking-wide mt-2 transition-all duration-500 group-hover:text-white">
                                        BỘ SƯU TẬP ĐỒNG HỒ
                                    </span>
                                </h1>
                            </div>

                            <!-- The Transformation -->
                            <div class="space-y-3 pt-4 border-t theme-border transition-all duration-500 group-hover:border-primary/50">
                                <p class="text-base sm:text-lg lg:text-xl text-muted-foreground leading-relaxed max-w-xl mx-auto lg:mx-0">
                                    Cảm nhận trọng lượng của thép tinh luyện trên cổ tay. Quan sát ánh chrome phản chiếu hoàn hảo dưới mọi góc độ. Mỗi chuyển động của kim là lời thì thầm của sự tinh tế.
                                </p>
                                <p class="text-base sm:text-lg lg:text-xl text-muted-foreground leading-relaxed max-w-xl mx-auto lg:mx-0 italic opacity-60 group-hover:opacity-100 transition-opacity">
                                    "Đồng hồ không chỉ đo thời gian. Nó đo sự xứng đáng của vật chủ."
                                </p>
                            </div>
                        </div>

                        <!-- Email Form -->
                        <div class="reveal reveal-delay-2 pt-4">
                            <p class="text-sm text-muted-foreground mb-4 uppercase tracking-[0.1em]">
                                Đăng ký nhận thông báo kích hoạt Drop
                            </p>
                            <waitlist-form></waitlist-form>
                        </div>
                    </div>

                    <!-- Watch Carousel -->
                    <div class="reveal reveal-delay-3 order-1 lg:order-2 relative group">
                        <div class="absolute inset-0 bg-primary/5 rounded-full blur-3xl transition-all duration-700 group-hover:bg-primary/10 group-hover:scale-110"></div>
                        <div class="relative transition-transform duration-500 group-hover:scale-[1.02]">
                            <watch-carousel></watch-carousel>
                        </div>
                    </div>
                </div>
            </div>
        </main>
`;
    }
  }
);
