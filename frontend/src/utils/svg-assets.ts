// MagTrade Vector Asset Library
// 纯代码生成的 SVG，用于替代缺失的商品图片

const COLORS = {
  stroke: '#333333',
  fill: '#1A1A1A',
  accent: '#FF3B30'
}

const PHONE_SVG = `
<svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect x="60" y="20" width="80" height="160" rx="12" fill="${COLORS.fill}" stroke="${COLORS.stroke}" stroke-width="2"/>
  <rect x="60" y="20" width="80" height="160" rx="12" stroke="white" stroke-opacity="0.1" stroke-width="2"/>
  <path d="M95 30H105" stroke="${COLORS.stroke}" stroke-width="2" stroke-linecap="round"/>
  <circle cx="100" cy="165" r="8" stroke="${COLORS.stroke}" stroke-width="2"/>
  <rect x="68" y="40" width="64" height="110" fill="#050505"/>
  <path d="M60 60L140 140" stroke="white" stroke-opacity="0.05" stroke-width="1"/>
</svg>`

const LAPTOP_SVG = `
<svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
  <path d="M30 50H170C172.209 50 174 51.7909 174 54V130H26V54C26 51.7909 27.7909 50 30 50Z" fill="${COLORS.fill}" stroke="${COLORS.stroke}" stroke-width="2"/>
  <rect x="34" y="58" width="132" height="64" fill="#050505"/>
  <path d="M20 130H180L190 150H10L20 130Z" fill="${COLORS.fill}" stroke="${COLORS.stroke}" stroke-width="2"/>
  <path d="M20 130H180" stroke="white" stroke-opacity="0.1" stroke-width="1"/>
</svg>`

const HEADPHONE_SVG = `
<svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
  <path d="M40 100V140C40 151.046 48.9543 160 60 160V160C71.0457 160 80 151.046 80 140V120H60C48.9543 120 40 111.046 40 100V100Z" fill="${COLORS.fill}" stroke="${COLORS.stroke}" stroke-width="2"/>
  <path d="M160 100V140C160 151.046 151.046 160 140 160V160C128.954 160 120 151.046 120 140V120H140C151.046 120 160 111.046 160 100V100Z" fill="${COLORS.fill}" stroke="${COLORS.stroke}" stroke-width="2"/>
  <path d="M40 100C40 66.8629 66.8629 40 100 40C133.137 40 160 66.8629 160 100" stroke="${COLORS.stroke}" stroke-width="2" stroke-linecap="round"/>
</svg>`

const BOX_SVG = `
<svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect x="50" y="50" width="100" height="100" fill="${COLORS.fill}" stroke="${COLORS.stroke}" stroke-width="2"/>
  <path d="M50 50L150 150" stroke="${COLORS.stroke}" stroke-width="1"/>
  <path d="M150 50L50 150" stroke="${COLORS.stroke}" stroke-width="1"/>
  <rect x="90" y="90" width="20" height="20" fill="${COLORS.accent}"/>
</svg>`

export const getSvgByType = (type: string): string => {
  let svg = BOX_SVG
  switch (type) {
    case 'phone': svg = PHONE_SVG; break;
    case 'laptop': svg = LAPTOP_SVG; break;
    case 'headphone': svg = HEADPHONE_SVG; break;
    case 'box': default: svg = BOX_SVG; break;
  }
  return `data:image/svg+xml;utf8,${encodeURIComponent(svg)}`
}

export const getSvgByProductName = (name: string): string => {
  const n = name.toLowerCase()
  let type = 'box'
  if (n.includes('phone') || n.includes('手机')) type = 'phone'
  else if (n.includes('mac') || n.includes('laptop') || n.includes('电脑')) type = 'laptop'
  else if (n.includes('pod') || n.includes('headphone') || n.includes('耳机')) type = 'headphone'
  
  return getSvgByType(type)
}