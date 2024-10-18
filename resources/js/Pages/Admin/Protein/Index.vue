<script setup>
import Pagination from '@/Components/Pagination.vue';
import Spinner from '@/Components/Spinner.vue';
import AuthenticatedLayout from '@/Layouts/AuthenticatedLayout.vue';
import { Head, Link, router, useForm } from '@inertiajs/vue3';
import { debounce } from 'lodash';
import { ref, watch } from 'vue';

const { proteins, search } = defineProps({
    proteins: Object,
    search: String
})

const importing = ref(false)
const processing = ref(false)

// Computed
const form = useForm({
    search: search
})

// Methods
function importArticles (id) {
    importing.value = true

    router.post(route('admin.articles.store'), { id }, {
        preserveScroll: true,
        preserveState: true,
        onSuccess: () => importing.value = false
    })
}

function processMutations (id) {
    processing.value = true

    router.post(route('admin.mutations.store'), { protein_id: id }, {
        preserveScroll: true,
        preserveState: true,
        onSuccess: () => processing.value = false
    })
}

// Watch
watch(
    () => form.search,
    debounce((newSearch, oldSearch) => {
        if (newSearch != oldSearch) {
            form.get(route('admin.proteins.index'), {
                preserveScroll: true,
                preserveState: true
            })
        }
    }, 500)
)
</script>

<template>
    <Head title="Dashboard" />

    <AuthenticatedLayout>
        <template #header>
            <h2 class="text-xl font-semibold leading-tight text-gray-800">
                Protein List
            </h2>
        </template>

        <div class="py-8">
            <div class="mx-auto max-w-7xl sm:px-6 lg:px-8">
                <div class="overflow-hidden bg-white shadow-sm sm:rounded-lg">
                    <div class="p-6 text-gray-900">
                        <!-- Protein List -->
                        <div class="pb-4 bg-white dark:bg-gray-900">
                            <label for="table-search" class="sr-only">Search</label>
                            <div class="relative mt-1">
                                <div class="absolute inset-y-0 rtl:inset-r-0 start-0 flex items-center ps-3 pointer-events-none">
                                    <svg class="w-4 h-4 text-gray-500 dark:text-gray-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20">
                                        <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"/>
                                    </svg>
                                </div>
                                <input v-model.trim="form.search" type="text" id="table-search" class="block pt-2 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg w-80 bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="Search for imported proteins">
                            </div>
                        </div>

                        <table class="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
                            <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                                <tr>
                                    <th scope="col" class="p-4" width="5%">
                                        id
                                    </th>
                                    <th scope="col" class="px-6 py-3" width="15%">
                                        Protein
                                    </th>
                                    <th scope="col" class="px-6 py-3" width="20%">
                                        PMIDS
                                    </th>
                                    <th scope="col" class="px-6 py-3" width="20%">
                                        Mutations
                                    </th>
                                    <th scope="col" class="px-6 py-3" width="20%">
                                        Action
                                    </th>
                                </tr>
                            </thead>
                            <tbody>
                                <template v-if="proteins?.data?.length">
                                    <tr v-for="protein in proteins.data" :key="protein.id" class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600">
                                        <td class="w-4 p-4">
                                            {{ protein.id }}
                                        </td>
                                        <th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                                            {{ protein.name }}
                                        </th>
                                        <td class="px-6 py-4">
                                            {{ protein.articles_count }}
                                        </td>
                                        <td class="px-6 py-4">
                                            {{ protein.mutations_count }}
                                        </td>
                                        <td class="px-6 py-4">
                                            <button v-if="!protein.articles_count" v-on:click.prevent="importArticles(protein.id)" class="rounded-md items-center bg-green-100 py-0.5 px-1 border border-transparent text-sm text-green-800 transition-all shadow-sm mr-2">
                                                <span v-if="importing" class="flex items-center">
                                                    <Spinner class="mr-2" />
                                                    Importing
                                                </span>

                                                <span v-else class="items-center">Import PMIDs</span>
                                            </button>

                                            <button v-on:click.prevent="processMutations(protein.id)" class="rounded-md items-center bg-red-100 py-0.5 px-1 border border-transparent text-sm text-red-800 transition-all shadow-sm">
                                                <span v-if="processing" class="flex items-center">
                                                    <Spinner class="mr-2" />
                                                    Processing
                                                </span>

                                                <span v-else class="items-center">Process Mutations</span>
                                            </button>
                                        </td>
                                    </tr>
                                </template>

                                <tr v-else class="text-center text-xl bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 underline">
                                    <td colspan="5" class="p-4">
                                        No Protein(s) Found.
                                    </td>
                                </tr>
                            </tbody>
                        </table>

                        <div v-if="proteins?.data?.length" class="mt-4">
                            <Pagination :pagination="proteins" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </AuthenticatedLayout>
</template>
