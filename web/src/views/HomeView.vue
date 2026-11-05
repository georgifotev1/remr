<script setup lang="ts">
import { useFetch } from '@vueuse/core'
import { Button } from 'primevue'

interface HealthResponse {
  status: string
  version: string
}

const { data, isFinished } = useFetch('/api/v1/health').json<HealthResponse>()
</script>

<template>
  <main>
    <h1>remr</h1>
    <p v-if="!isFinished">Loading...</p>
    <p v-else>API: Status: {{ data?.status }} Version: {{ data?.version }}</p>
    <Button>Hello</Button>
  </main>
</template>
