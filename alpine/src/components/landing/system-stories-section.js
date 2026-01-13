// System Stories Section Component for Alpine.js
document.addEventListener("alpine:init", () => {
  Alpine.data("systemStories", () => ({
    activeStory: 0,
    stories: [
      {
        id: "SHARP",
        title: "SẮC LẸM",
        subtitle: "The Edge",
        description: "Đường nét tương lai. Không chi tiết thừa.",
        stat: "PERFECTION",
        color: "text-white",
      },
      {
        id: "GLOW",
        title: "BÓNG BẨY",
        subtitle: "The Aura",
        description: "Ánh sáng quyền lực. Thu hút mọi ánh nhìn.",
        stat: "DOMINANCE",
        color: "text-accent",
      },
      {
        id: "SILENT",
        title: "TĨNH LẶNG",
        subtitle: "The Soul",
        description: "Sang trọng tối thượng. Đẳng cấp thầm lặng.",
        stat: "ELEGANCE",
        color: "text-zinc-400",
      },
      {
        id: "BOLD",
        title: "TÁO BẠO",
        subtitle: "The Vision",
        description: "Phá vỡ quy tắc. Biểu tượng mới.",
        stat: "REBELLION",
        color: "text-primary",
      },
    ],
  }));
});

customElements.define(
  "system-stories-section",
  class extends HTMLElement {
    connectedCallback() {
      this.innerHTML = /*html*/ `
        <section x-data="systemStories" class="relative py-32 px-6 overflow-hidden bg-black">
            <!-- Background Ambient Glow -->
            <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-primary/5 rounded-full blur-[160px] pointer-events-none"></div>
            
            <div class="max-w-7xl mx-auto relative z-10">
                <div class="flex flex-col lg:flex-row gap-20 items-start">
                    <!-- Left: Story Navigation -->
                    <div class="w-full lg:w-1/3 space-y-12 reveal">
                        <div class="space-y-4">
                            <h2 class="text-5xl font-black tracking-tighter text-white leading-none">
                                LƯU TRỮ<br/><span class="text-primary">PHONG CÁCH</span>
                            </h2>
                            <div class="h-1 w-20 bg-primary"></div>
                            <p class="text-zinc-500 font-mono text-xs uppercase tracking-[0.3em]">Trạng thái: Đang tải thông số thiết kế...</p>
                        </div>

                        <div class="flex flex-col gap-4">
                            <template x-for="(story, index) in stories" :key="story.id">
                                <button 
                                    @click="activeStory = index"
                                    class="group relative text-left py-4 transition-all duration-500 border-b border-white/5"
                                    :class="activeStory === index ? 'pl-8' : 'hover:pl-4 opacity-40 hover:opacity-100'"
                                >
                                    <div 
                                        class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-0 bg-primary transition-all duration-500"
                                        :class="activeStory === index ? 'h-2/3' : 'group-hover:h-1/3'"
                                    ></div>
                                    <span class="block text-[10px] font-mono text-zinc-500 mb-1" x-text="story.id"></span>
                                    <span 
                                        class="text-2xl font-black tracking-tight uppercase transition-colors"
                                        :class="activeStory === index ? 'text-white' : 'text-zinc-400'"
                                        x-text="story.title"
                                    ></span>
                                </button>
                            </template>
                        </div>
                    </div>

                    <!-- Right: Immersive Content -->
                    <div class="w-full lg:w-2/3 min-h-[500px] flex items-center justify-center relative">
                        <template x-for="(story, index) in stories" :key="story.id">
                            <div 
                                x-show="activeStory === index"
                                x-transition:enter="transition ease-out duration-700 delay-300"
                                x-transition:enter-start="opacity-0 translate-y-12 scale-95"
                                x-transition:enter-end="opacity-100 translate-y-0 scale-100"
                                x-transition:leave="transition ease-in duration-300 absolute"
                                x-transition:leave-start="opacity-100 scale-100"
                                x-transition:leave-end="opacity-0 scale-110"
                                class="w-full space-y-12"
                            >
                                <div class="relative">
                                    <span class="absolute -top-20 -left-10 text-[12rem] font-black text-white/[0.02] select-none pointer-events-none" x-text="story.id"></span>
                                    <div class="space-y-6 relative z-10">
                                        <p class="text-primary font-mono text-sm tracking-[0.4em]" x-text="story.subtitle"></p>
                                        <h3 class="text-4xl sm:text-6xl font-black text-white uppercase leading-tight max-w-2xl" x-text="story.description"></h3>
                                    </div>
                                </div>

                                <div class="flex items-center gap-12 pt-12 border-t border-white/10">
                                    <div class="space-y-1">
                                        <p class="text-zinc-500 text-[10px] font-mono uppercase tracking-widest">Metric</p>
                                        <p class="text-2xl font-black text-white" x-text="story.stat"></p>
                                    </div>
                                    <div class="h-12 w-px bg-white/10"></div>
                                    <div class="space-y-1">
                                        <p class="text-zinc-500 text-[10px] font-mono uppercase tracking-widest">Status</p>
                                        <p class="text-2xl font-black text-primary animate-pulse">ACTIVE</p>
                                    </div>
                                    
                                    <div class="ml-auto">
                                        <div class="w-16 h-16 rounded-full border border-primary/30 flex items-center justify-center group cursor-pointer hover:bg-primary transition-all duration-500">
                                            <svg class="w-6 h-6 text-primary group-hover:text-black transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l7 7m-7-7v18"></path>
                                            </svg>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </template>
                    </div>
                </div>
            </div>

            <!-- Decorative Elements -->
            <div class="absolute bottom-0 left-0 w-full h-px bg-gradient-to-r from-transparent via-white/10 to-transparent"></div>
        </section>
      `;
    }
  }
);
