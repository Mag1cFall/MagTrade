<template>
  <div class="testimonials-container">
    <div class="cards-wrapper" ref="wrapper">
      <div v-for="(item, i) in testimonials" :key="i" class="testimonial-card">
        <div class="quote-symbol">"</div>
        <p class="feedback-text">{{ item.text }}</p>
        <div class="user-info">
          <div class="avatar">{{ item.name.charAt(0) }}</div>
          <div>
            <div class="name">{{ item.name }}</div>
            <div class="role">{{ item.role }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const testimonials = [
  { name: 'Elon M.', role: 'Tech Entrepreneur', text: 'The latency on MagTrade is absolutely unreal. It feels like direct neural interface trading.' },
  { name: 'Sarah C.', role: 'Pro Scalper', text: "I've used every platform out there. Nothing comes close to the execution speed here." },
  { name: 'David L.', role: 'Hedge Fund Mgr', text: 'Secure, reliable, and ridiculously fast. The new standard for high-frequency systems.' },
  { name: 'Jessica W.', role: 'Crypto Whale', text: 'Finally, a platform that can handle the volatility without crashing.' },
  { name: 'Wei Z.', role: 'Quant Analyst', text: "The Obsidian engine matching efficiency is the best I've seen. Zero slippage." },
  { name: 'Akira T.', role: 'System Architect', text: 'From a technical perspective, their distributed node layout is a work of art.' },
];
</script>

<style scoped>
.testimonials-container {
  overflow-x: auto;
  padding-bottom: 20px;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
}
.testimonials-container::-webkit-scrollbar {
  display: none;
}

.cards-wrapper {
  display: flex;
  gap: 24px;
  width: max-content;
  padding: 0 20px;
}

.testimonial-card {
  width: 320px;
  background: #111;
  border: 1px solid #222;
  padding: 30px;
  border-radius: 16px;
  position: relative;
  flex-shrink: 0;
  transition: transform 0.3s;
}

.testimonial-card:hover {
  transform: translateY(-5px);
  border-color: #333;
}

.quote-symbol {
  font-size: 60px;
  color: #222;
  font-family: serif;
  line-height: 1;
  position: absolute;
  top: 10px;
  right: 20px;
}

.feedback-text {
  font-size: 14px;
  color: #ccc;
  line-height: 1.6;
  margin-bottom: 24px;
  position: relative;
  z-index: 1;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 40px;
  height: 40px;
  background: #e33535;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
}

.name {
  color: #fff;
  font-weight: bold;
  font-size: 14px;
}

.role {
  color: #666;
  font-size: 12px;
}
</style>
