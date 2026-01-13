// Cinematic Video Section Component for Alpine.js
document.addEventListener("alpine:init", () => {
  Alpine.data("cinematicVideoSection", () => ({
    isMuted: true,
    videoScale: 1,

    toggleMute() {
      this.isMuted = !this.isMuted;
      this.$refs.video.muted = this.isMuted;
    },

    init() {
      window.addEventListener("scroll", () => {
        const rect = this.$el.getBoundingClientRect();
        const scrollPercent = Math.max(
          0,
          Math.min(
            1,
            (window.innerHeight - rect.top) / (window.innerHeight + rect.height)
          )
        );
        this.videoScale = 1 + scrollPercent * 0.1;
      });
    },
  }));
});

customElements.define(
  "cinematic-video-section",
  class extends HTMLElement {
    connectedCallback() {
      this.innerHTML = /*html*/ `
    <section x-data="cinematicVideoSection" class="relative overflow-hidden">
        <!-- Video Section -->
        <div class="relative min-h-[50vh] sm:min-h-[60vh] lg:min-h-[80vh] overflow-hidden">
            <div class="absolute inset-0 transition-transform duration-300 ease-out" :style="'transform: scale(' + videoScale + ')'">
                <video
                    x-ref="video"
                    class="h-full w-full object-cover"
                    src="/videos/logo.mp4"
                    autoplay
                    muted
                    loop
                    playsinline
                >
                    Your browser does not support the video tag.
                </video>
            </div>
            
            <!-- Mute Toggle -->
            <button @click="toggleMute" class="absolute bottom-8 right-8 z-20 p-4 rounded-full bg-black/40 backdrop-blur-md border border-white/10 text-white hover:bg-white/20 transition-all group">
                <template x-if="isMuted">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z"></path>
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2"></path>
                    </svg>
                </template>
                <template x-if="!isMuted">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 12.728M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z"></path>
                    </svg>
                </template>
            </button>
        </div>

        <!-- Text Content Below Video -->
        <div class="relative bg-background py-16 sm:py-24 lg:py-32 px-4 sm:px-6 md:px-8 lg:px-12">
            <div class="max-w-5xl mx-auto">
                <!-- Collection Announcement & Countdown -->
                <div class="reveal text-center mb-24">
                    <div class="inline-flex items-center gap-3 px-4 py-1.5 rounded-full bg-white/[0.03] border border-white/10 backdrop-blur-md mb-8">
                        <span class="relative flex h-2 w-2">
                            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                            <span class="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
                        </span>
                        <span class="text-[10px] font-mono text-white/60 uppercase tracking-[0.2em]">Limited Drop 01</span>
                    </div>

                    <h2 class="text-4xl sm:text-5xl lg:text-6xl font-black text-white uppercase tracking-tighter mb-12 relative group">
                        <span class="relative z-10" x-show="!isDropLive">SẮP <span class="text-primary">KÍCH HOẠT</span></span>
                        <span class="relative z-10" x-show="isDropLive">ĐÃ <span class="text-green-500">KÍCH HOẠT</span></span>
                        <span class="absolute inset-0 text-white/5 blur-xl group-hover:text-primary/20 transition-colors duration-700 select-none" x-show="!isDropLive">SẮP KÍCH HOẠT</span>
                        <span class="absolute inset-0 text-green-500/5 blur-xl group-hover:text-green-500/20 transition-colors duration-700 select-none" x-show="isDropLive">ĐÃ KÍCH HOẠT</span>
                    </h2>
                    
                    <!-- Elegant Dropdown Clock -->
                    <div class="flex flex-col items-center space-y-12">
                        <div class="space-y-2">
                            <p class="text-[10px] sm:text-xs text-primary uppercase tracking-[0.5em] font-bold animate-pulse">
                                Giao thức: Drop-Sequence-01
                            </p>
                            <p class="text-xs text-muted-foreground uppercase tracking-[0.2em]" x-show="!isDropLive">
                                Đang khởi tạo căn chỉnh thời gian
                            </p>
                            <p class="text-xs text-green-500 uppercase tracking-[0.2em] font-bold animate-pulse" x-show="isDropLive">
                                ĐÃ KÍCH HOẠT
                            </p>
                        </div>
                        
                        <div class="flex items-center justify-center gap-3 sm:gap-6" x-show="!isDropLive">
                            <template x-for="(val, unit) in countdown" :key="unit">
                                <div class="flex flex-col items-center group reveal clock-unit" 
                                     :style="'animation-delay: ' + (unit === 'd' ? '0.1s' : unit === 'h' ? '0.2s' : unit === 'm' ? '0.3s' : '0.4s')"
                                     :class="unit === 'd' ? 'reveal-delay-1' : unit === 'h' ? 'reveal-delay-2' : unit === 'm' ? 'reveal-delay-3' : 'reveal-delay-4'">
                                    <div class="relative overflow-hidden rounded-2xl bg-gradient-to-b from-white/[0.05] to-transparent border border-white/10 p-5 sm:p-8 min-w-[80px] sm:min-w-[120px] transition-all duration-700 group-hover:border-primary/40 group-hover:shadow-[0_0_30px_rgba(59,130,246,0.1)] group-hover:-translate-y-3">
                                        <!-- Dropdown Number -->
                                        <div class="relative z-10">
                                            <span class="text-4xl sm:text-6xl font-black text-white tracking-tighter font-mono leading-none" x-text="val"></span>
                                        </div>
                                        
                                        <!-- Glass Reflection -->
                                        <div class="absolute inset-0 bg-gradient-to-tr from-transparent via-white/[0.02] to-transparent pointer-events-none"></div>
                                        
                                        <!-- Top Slot Line -->
                                        <div class="absolute top-0 left-1/2 -translate-x-1/2 w-1/2 h-px bg-gradient-to-r from-transparent via-primary/50 to-transparent"></div>
                                    </div>
                                    <span class="mt-4 text-[9px] sm:text-[10px] font-mono text-zinc-500 uppercase tracking-[0.3em] group-hover:text-primary transition-colors" x-text="unit === 'd' ? 'Days' : unit === 'h' ? 'Hours' : unit === 'm' ? 'Mins' : 'Secs'"></span>
                                </div>
                            </template>
                        </div>

                        <div class="flex flex-col items-center gap-6 pt-8">
                            <div class="flex items-center gap-4">
                                <div class="h-px w-16 bg-gradient-to-r from-transparent to-white/20"></div>
                                <span class="text-[10px] font-mono text-zinc-500 uppercase tracking-[0.5em]">Awaiting Activation</span>
                                <div class="h-px w-16 bg-gradient-to-l from-transparent to-white/20"></div>
                            </div>
                            
                            <!-- Animated Progress Bar -->
                            <div class="w-48 h-1 bg-white/5 rounded-full overflow-hidden relative">
                                <div class="absolute inset-0 bg-primary/40 animate-[shimmer_2s_infinite] w-1/2 rounded-full"></div>
                            </div>
                        </div>
                        </div>
                    </div>
                </div>

                <!-- Story Hook - Opening -->
                <div class="text-foreground mb-12 space-y-8">
                    <div class="reveal reveal-delay-1">
                        <p class="text-sm sm:text-base lg:text-lg text-muted-foreground uppercase tracking-[0.2em] mb-4">
                            Trong thế giới của những điều tầm thường
                        </p>
                        <h1 class="font-display text-4xl sm:text-5xl lg:text-6xl xl:text-7xl leading-tight font-black tracking-[-0.02em] text-foreground">
                            KHÔNG CHỈ LÀ VẬT THỂ.
                        </h1>
                    </div>
                    
                    <div class="reveal reveal-delay-2">
                        <h1 class="font-display text-4xl sm:text-5xl lg:text-6xl xl:text-7xl leading-tight font-black tracking-[-0.02em] text-foreground">
                            ĐÂY LÀ MỘT TRẢI NGHIỆM.
                        </h1>
                    </div>

                    <div class="reveal reveal-delay-3 max-w-2xl">
                        <p class="text-lg sm:text-xl lg:text-2xl text-muted-foreground leading-relaxed">
                            Mỗi chiếc đồng hồ là một hành trình khám phá, một câu chuyện được kể qua từng chuyển động của kim giờ.
                            Sự tinh tế không nằm ở bề ngoài, mà ở cảm giác mà nó mang lại cho người sở hữu.
                        </p>
                    </div>
                </div>

                <!-- Scroll Indicator -->
                <div class="reveal reveal-delay-4 flex items-center justify-center gap-4 text-muted-foreground pt-12">
                    <span class="text-xs uppercase tracking-[0.4em] font-bold">Scroll to explore</span>
                    <div class="w-px h-12 bg-gradient-to-b from-primary to-transparent animate-bounce"></div>
                </div>
            </div>
        </div>
    </section>
`;
    }
  }
);
