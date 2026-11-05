import { definePreset } from '@primeuix/themes'
import Aura from '@primeuix/themes/aura'

export const RemyPreset = definePreset(Aura, {
  primitive: {
    borderRadius: {
      none: '0',
      xs: '4px',
      sm: '6px',
      md: '8px',
      lg: '10px',
      xl: '14px',
    },
  },
  semantic: {
    primary: {
      50: '#fff4ed',
      100: '#ffe7d7',
      200: '#ffcdae',
      300: '#ffaa7a',
      400: '#ff823f',
      500: '#ed5700',
      600: '#cb3a00',
      700: '#a42f00',
      800: '#7e2401',
      900: '#612009',
      950: '#3a1003',
      color: '{zinc.900}',
      contrastColor: '#ffffff',
      hoverColor: '{zinc.800}',
      activeColor: '{zinc.700}',
    },
    surface: {
      0: '#ffffff',
      50: '{zinc.50}',
      100: '{zinc.100}',
      200: '{zinc.200}',
      300: '{zinc.300}',
      400: '{zinc.400}',
      500: '{zinc.500}',
      600: '{zinc.600}',
      700: '{zinc.700}',
      800: '{zinc.800}',
      900: '{zinc.900}',
      950: '{zinc.950}',
    },
    // Orange focus ring, like the prototype's --ring.
    focusRing: {
      width: '2px',
      style: 'solid',
      color: 'color-mix(in srgb, {primary.500}, transparent 40%)',
      offset: '2px',
      shadow: 'none',
    },
    formField: {
      paddingX: '0.75rem',
      paddingY: '0.5rem',
      borderColor: '{surface.300}',
      hoverBorderColor: '{surface.400}',
      focusBorderColor: '{primary.500}',
      color: '{surface.950}',
      placeholderColor: '{surface.500}',
    },
    text: {
      color: '{surface.950}',
      hoverColor: '{surface.950}',
      mutedColor: '{surface.600}',
      hoverMutedColor: '{surface.700}',
    },
    content: {
      borderRadius: '{border.radius.lg}',
      borderColor: '{surface.200}',
    },
    highlight: {
      background: '{primary.50}',
      focusBackground: '{primary.100}',
      color: '{primary.700}',
      focusColor: '{primary.800}',
    },
  },
})
