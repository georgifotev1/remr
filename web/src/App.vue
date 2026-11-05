<script setup lang="ts">
import { ConfigProvider } from 'reka-ui'
import { RouterView } from 'vue-router'
</script>

<template>
  <ConfigProvider locale="bg">
    <RouterView />
  </ConfigProvider>
</template>
