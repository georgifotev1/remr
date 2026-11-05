import PrimeVue from 'primevue/config'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/main.css'
import { bg } from './theme/locale-bg'
import { RemyPreset } from './theme/preset'

const app = createApp(App)

app.use(PrimeVue, {
  locale: bg,
  theme: {
    preset: RemyPreset,
    options: {
      darkModeSelector: '.app-dark',
    },
  },
})
app.use(router)
app.mount('#app')
