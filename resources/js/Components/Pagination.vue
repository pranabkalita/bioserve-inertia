<script setup>
    import { Link } from '@inertiajs/vue3'

    defineProps({
        pagination: Object
    })
</script>

<template>
    <div class="flex justify-between">
        <p class="flex items-center justify-center px-3 py-2 text-sm rounded-lg text-gray-600">Total: {{ pagination.total }}</p>

        <nav class="relative flex justify-center">
            <div v-for="link in pagination.links" :key="link.label">
                <Link
                    preserve-scroll
                    :href="link.url ?? ''"
                    v-html="link.label"
                    class="flex items-center justify-center px-3 py-2 text-sm rounded-lg text-gray-600"
                    :class="{ 'bg-gray-200': link.active, '!text-gray-300': !link.url }"
                />
            </div>
        </nav>
    </div>
</template>
