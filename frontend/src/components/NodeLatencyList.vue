<template>
  <div class="node-latency-list bg-black/40 backdrop-blur border border-white/10 rounded-xl p-6">
    <h3 class="text-xs font-bold text-gray-500 uppercase tracking-widest mb-4">Network Status</h3>
    
    <div class="space-y-3">
      <div v-for="node in nodes" :key="node.id" class="flex items-center justify-between text-sm group cursor-pointer hover:bg-white/5 p-2 rounded transition-colors">
        <div class="flex items-center gap-3">
          <div 
            class="w-2 h-2 rounded-full"
            :class="node.status === 'operational' ? 'bg-green-500 animate-pulse' : 'bg-yellow-500'"
          ></div>
          <span class="font-mono text-gray-300 group-hover:text-white">{{ node.name }}</span>
        </div>
        
        <div class="flex items-center gap-4">
          <span class="text-xs text-gray-600 font-mono">{{ node.region }}</span>
          <span 
            class="font-mono font-bold w-12 text-right"
            :class="getLatencyColor(node.latency)"
          >{{ node.latency }}ms</span>
        </div>
      </div>
    </div>

    <div class="mt-6 pt-4 border-t border-white/5 flex justify-between items-center text-xs">
      <span class="text-gray-500">Global Uptime: 99.99%</span>
      <span class="text-green-500 flex items-center gap-1">
        <div class="w-1.5 h-1.5 bg-current rounded-full"></div> All Systems Operational
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

const nodes = ref([
  { id: 'tokyo', name: 'Tokyo-01', region: 'AP-NE', status: 'operational', latency: 45 },
  { id: 'sg', name: 'Singapore-X', region: 'AP-SE', status: 'operational', latency: 82 },
  { id: 'frankfurt', name: 'Frankfurt-M', region: 'EU-CE', status: 'operational', latency: 140 },
  { id: 'ny', name: 'NewYork-Core', region: 'US-EA', status: 'operational', latency: 180 },
]);

// 模拟延迟波动
let timer: any;

const updateLatency = () => {
  nodes.value.forEach(node => {
    // 随机波动 +/- 5ms
    const noise = Math.floor(Math.random() * 11) - 5;
    node.latency = Math.max(5, node.latency + noise);
  });
};

const getLatencyColor = (ms: number) => {
  if (ms < 50) return 'text-green-400';
  if (ms < 100) return 'text-yellow-400';
  return 'text-red-400';
};

onMounted(() => {
  timer = setInterval(updateLatency, 2000);
});

onUnmounted(() => {
  clearInterval(timer);
});
</script>
