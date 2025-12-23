<template>
  <div class="relative w-full h-[400px] md:h-[600px] overflow-hidden bg-[#050505] flex items-center justify-center border border-white/5 rounded-xl group">
    <!-- Grid Background -->
    <div class="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-[size:40px_40px] [mask-image:radial-gradient(ellipse_at_center,black_40%,transparent_80%)]"></div>

    <!-- Map Dots -->
    <svg viewBox="0 0 800 400" class="w-full h-full opacity-60">
      <g v-for="(dot, i) in dots" :key="i" :transform="`translate(${dot.x}, ${dot.y})`">
        <circle r="1.5" fill="#333" />
      </g>
      
      <!-- Active Nodes -->
      <g v-for="(node, i) in nodes" :key="`node-${i}`" :transform="`translate(${node.x}, ${node.y})`">
        <circle r="3" fill="#e33535" class="animate-pulse" />
        <circle r="8" fill="none" stroke="#e33535" stroke-opacity="0.3" class="node-ring" />
        <circle r="16" fill="none" stroke="#e33535" stroke-opacity="0.1" class="node-ring-lg" />
        
        <!-- Connection Lines -->
        <line
          v-if="i < nodes.length - 1"
          :x1="0" :y1="0"
          :x2="(nodes[i+1]?.x || 0) - node.x" :y2="(nodes[i+1]?.y || 0) - node.y"
          stroke="#e33535"
          stroke-width="0.5"
          stroke-opacity="0.2"
        />
      </g>
    </svg>

    <!-- Overlay UI -->
    <div class="absolute bottom-8 left-8 flex flex-col gap-2">
      <div class="flex items-center gap-2">
        <span class="w-2 h-2 bg-accent rounded-full animate-ping"></span>
        <span class="text-xs font-mono text-accent uppercase tracking-widest">Live Network Status</span>
      </div>
      <div class="text-3xl font-bold text-white tracking-tighter">GLOBAL <span class="text-secondary">NODES</span></div>
      <div class="text-sm text-secondary font-mono">
        Active Regions: <span class="text-white">14</span> | Latency: <span class="text-accent">< 35ms</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const dots = Array.from({ length: 200 }, () => ({
  x: Math.random() * 800,
  y: Math.random() * 400
}))

const nodes = [
  { x: 200, y: 150, name: 'NYC' },
  { x: 250, y: 180, name: 'LDN' },
  { x: 400, y: 120, name: 'FRA' },
  { x: 600, y: 160, name: 'SGP' },
  { x: 650, y: 140, name: 'TYO' },
  { x: 150, y: 200, name: 'SFO' }
]
</script>

<style scoped>
.node-ring {
  animation: ripple 2s cubic-bezier(0, 0.2, 0.8, 1) infinite;
}
.node-ring-lg {
  animation: ripple 2s cubic-bezier(0, 0.2, 0.8, 1) infinite;
  animation-delay: 0.5s;
}

@keyframes ripple {
  0% { transform: scale(0.5); opacity: 1; }
  100% { transform: scale(1.5); opacity: 0; }
}
</style>