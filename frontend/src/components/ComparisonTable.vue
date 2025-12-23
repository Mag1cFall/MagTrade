<template>
  <div class="w-full">
    <div class="w-full grid grid-cols-1 md:grid-cols-4 gap-1 text-sm bg-surface border border-white/10 rounded-lg overflow-hidden">
      <!-- Headers -->
      <div class="p-4 md:p-6 bg-surface-light/50 border-b border-white/10 text-secondary font-bold uppercase tracking-widest flex items-end">
        Architecture
      </div>
      <div class="hidden md:flex p-6 bg-surface-light/50 border-b border-white/10 text-tertiary font-bold uppercase tracking-widest text-center flex-col justify-end opacity-50">
        <div class="text-xs mb-1">Legacy</div>
        <div>Monolith</div>
      </div>
      <div class="hidden md:flex p-6 bg-surface-light/50 border-b border-white/10 text-tertiary font-bold uppercase tracking-widest text-center flex-col justify-end opacity-70">
        <div class="text-xs mb-1">Standard</div>
        <div>Microservices</div>
      </div>
      <div class="relative p-4 md:p-6 bg-surface-light border-b border-accent text-white font-bold uppercase tracking-widest text-center flex flex-col justify-end">
        <div class="absolute top-0 left-0 w-full h-1 bg-accent shadow-[0_0_15px_#e33535]"></div>
        <div class="text-xs text-accent mb-1">MagTrade</div>
        <div>Obsidian Engine</div>
      </div>

      <!-- Row 1 -->
      <div class="p-4 px-6 border-b border-white/5 text-secondary flex items-center">
        Transaction Throughput
      </div>
      <div class="hidden md:flex p-4 px-6 border-b border-white/5 text-tertiary text-center items-center justify-center">
        ~100 TPS
      </div>
      <div class="hidden md:flex p-4 px-6 border-b border-white/5 text-tertiary text-center items-center justify-center">
        ~1,000 TPS
      </div>
      <div class="p-4 px-6 border-b border-white/5 bg-accent/5 text-white font-mono font-bold text-center flex items-center justify-center border-l border-r border-accent/20">
        50,000+ TPS
      </div>

      <!-- Row 2 -->
      <div class="p-4 px-6 border-b border-white/5 text-secondary flex items-center">
        Network Latency
      </div>
      <div class="hidden md:flex p-4 px-6 border-b border-white/5 text-tertiary text-center items-center justify-center">
        200ms+
      </div>
      <div class="hidden md:flex p-4 px-6 border-b border-white/5 text-tertiary text-center items-center justify-center">
        50ms+
      </div>
      <div class="p-4 px-6 border-b border-white/5 bg-accent/5 text-white font-mono font-bold text-center flex items-center justify-center border-l border-r border-accent/20">
        < 5ms
      </div>

      <!-- Row 3 -->
      <div class="p-4 px-6 border-b border-white/5 text-secondary flex items-center">
        Conflict Resolution
      </div>
      <div class="hidden md:flex p-4 px-6 border-b border-white/5 text-tertiary text-center items-center justify-center">
        DB Locks
      </div>
      <div class="hidden md:flex p-4 px-6 border-b border-white/5 text-tertiary text-center items-center justify-center">
        Optimistic
      </div>
      <div class="p-4 px-6 border-b border-white/5 bg-accent/5 text-white font-bold text-center flex items-center justify-center border-l border-r border-accent/20">
        Atomic Lua
      </div>

      <!-- Row 4 -->
      <div class="p-4 px-6 text-secondary flex items-center">
        AI Tactical Support
      </div>
      <div class="hidden md:flex p-4 px-6 text-tertiary text-center items-center justify-center">
        None
      </div>
      <div class="hidden md:flex p-4 px-6 text-tertiary text-center items-center justify-center">
        Basic Stats
      </div>
      <div class="p-4 px-6 bg-accent/5 text-accent font-bold text-center flex items-center justify-center border-l border-r border-accent/20">
        DeepSeek Neural
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>