<template>
  <div class="tech-stack-grid">
    <div v-for="(tech, i) in techStack" :key="i" class="tech-item group">
      <div class="icon-box">
        <component :is="tech.icon" class="w-8 h-8 text-white group-hover:text-accent transition-colors" />
      </div>
      <div class="tech-info">
        <div class="tech-name">{{ tech.name }}</div>
        <div class="tech-desc">{{ tech.desc }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Cpu, Database, Globe, Network, Shield, Zap, Cloud, Code } from 'lucide-vue-next';

const techStack = [
  { name: 'Golang', desc: 'High-performance Backend', icon: Zap },
  { name: 'PostgreSQL', desc: 'Relational Database', icon: Database },
  { name: 'Redis', desc: 'In-memory Cache', icon: Cpu },
  { name: 'Kafka', desc: 'Message Queue', icon: Network },
  { name: 'Vue 3', desc: 'Reactive Frontend', icon: Code },
  { name: 'Docker', desc: 'Containerization', icon: Cloud },
  { name: 'WebSocket', desc: 'Real-time Transport', icon: Globe },
  { name: 'JWT', desc: 'Secure Auth', icon: Shield },
];
</script>

<style scoped>
.tech-stack-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

@media (min-width: 768px) {
  .tech-stack-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

.tech-item {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  padding: 24px;
  border-radius: 12px;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 16px;
}

.tech-item:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-5px);
}

.icon-box {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.05);
  transition: border-color 0.3s;
}

.tech-item:hover .icon-box {
  border-color: #e33535;
}

.tech-name {
  font-weight: bold;
  color: #fff;
  font-size: 16px;
}

.tech-desc {
  font-size: 12px;
  color: #888;
}
</style>
