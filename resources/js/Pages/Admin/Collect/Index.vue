<script setup>
import DangerButton from '@/Components/DangerButton.vue';
import InputError from '@/Components/InputError.vue';
import Pagination from '@/Components/Pagination.vue';
import PrimaryButton from '@/Components/PrimaryButton.vue';
import AuthenticatedLayout from '@/Layouts/AuthenticatedLayout.vue';
import { Head, router, useForm } from '@inertiajs/vue3';
import { pickBy } from 'lodash';
import { onMounted, reactive, ref } from 'vue';

const { pmidCount, search, proteins } = defineProps({
    search: String,
    pmidCount: Number,
    proteins: Object
})

const processing = ref(false)
const importing = ref(false)

// Computed
const params = reactive({
    search: ''
})

// Mounted
onMounted(() => {
    params.search = search
})

function searchProtein () {
    processing.value = true

    router.get(route('collect.index'), pickBy(params), {
        preserveState: true,
        onFinish: () => {
            processing.value = false
        }
    })
}

function importProteinData () {
    importing.value = true

    router.post(route('protein.store'), params, {
        preserveState: true,
        onFinish: () => {
            importing.value = false
        }
    })
}
</script>

<template>
    <Head title="Dashboard" />

    <AuthenticatedLayout>
        <template #header>
            <h2 class="text-xl font-semibold leading-tight text-gray-800">
                Collect Data
            </h2>
        </template>

        <div class="py-8">
            <div class="mx-auto max-w-7xl sm:px-6 lg:px-8">
                <div class="overflow-hidden bg-white shadow-sm sm:rounded-lg">
                    <div class="p-6 text-gray-900">
                        <!-- Protein Search -->
                        <form v-on:submit.prevent="searchProtein">
                            <div class="flex items-center w-full mx-auto">
                                <label for="simple-search" class="sr-only">Search</label>
                                <div class="relative w-full">
                                    <div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
                                        <svg class="w-4 h-4 text-gray-500 dark:text-gray-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 18 20">
                                            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5v10M3 5a2 2 0 1 0 0-4 2 2 0 0 0 0 4Zm0 10a2 2 0 1 0 0 4 2 2 0 0 0 0-4Zm12 0a2 2 0 1 0 0 4 2 2 0 0 0 0-4Zm0 0V6a3 3 0 0 0-3-3H9m1.5-2-2 2 2 2" />
                                        </svg>
                                    </div>

                                    <input v-model.trim="params.search" id="search" type="text" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full ps-10 p-2.5  dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="Search protein. Eg. BRCA1" required />
                                </div>

                                <button :class="{ 'opacity-25 cursor-wait': processing || importing }" :disabled="processing || importing" type="submit" class="p-2.5 ms-2 text-sm font-medium text-white bg-blue-700 rounded-lg border border-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
                                    <svg v-if="!processing && !importing" class="w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20">
                                        <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z" />
                                    </svg>

                                    <svg v-else class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24">
                                        <path fill="currentColor" d="M12,4a8,8,0,0,1,7.89,6.7A1.53,1.53,0,0,0,21.38,12h0a1.5,1.5,0,0,0,1.48-1.75,11,11,0,0,0-21.72,0A1.5,1.5,0,0,0,2.62,12h0a1.53,1.53,0,0,0,1.49-1.3A8,8,0,0,1,12,4Z">
                                            <animateTransform attributeName="transform" dur="0.75s" repeatCount="indefinite" type="rotate" values="0 12 12;360 12 12" />
                                        </path>
                                    </svg>
                                </button>
                            </div>

                            <InputError class="mt-2" :message="$page.props.errors.search" />

                            <div v-if="pmidCount && !processing" class="mt-2 flex">
                                <p class="text-green-600 dark:text-green-500">Found {{ pmidCount }} PMIDs for Protein : {{ params.search }}</p>

                                <button v-on:click.prevent="importProteinData" :disabled="importing" :class="{ 'opacity-25 cursor-wait': importing }"  type="button" class="flex text-white bg-blue-700 hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 font-medium rounded-full text-sm px-5 py-0.5 text-center ml-2 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800 items-center">
                                    <template v-if="importing">
                                        <svg v-if="importing" class="w-3 h-3 mr-2" xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24">
                                            <path fill="currentColor" d="M12,4a8,8,0,0,1,7.89,6.7A1.53,1.53,0,0,0,21.38,12h0a1.5,1.5,0,0,0,1.48-1.75,11,11,0,0,0-21.72,0A1.5,1.5,0,0,0,2.62,12h0a1.53,1.53,0,0,0,1.49-1.3A8,8,0,0,1,12,4Z">
                                                <animateTransform attributeName="transform" dur="0.75s" repeatCount="indefinite" type="rotate" values="0 12 12;360 12 12" />
                                            </path>
                                        </svg>
                                    <span>Importing ...</span>
                                    </template>

                                    <span v-else>Import</span>
                                </button>

                                <button :disabled="importing" :class="{ 'opacity-25 cursor-wait': importing }" type="button" class="text-white bg-red-700 hover:bg-red-800 focus:outline-none focus:ring-4 focus:ring-red-300 font-medium rounded-full text-sm px-5 py-0.5 text-center me-2 mb-2 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-red-900">Cancel</button>

                            </div>

                            <div v-if="!pmidCount && !processing && search" class="mt-2 text-sm text-red-600 dark:text-red-500">
                                No PMIDs found for searched protein. Try another protein.
                            </div>
                        </form>

                        <!-- Protein List -->
                        <template v-if="proteins?.data?.length">
                            <h1 class="my-5 text-4xl font-extrabold text-center text-gray-900 dark:text-white md:text-4xl lg:text-4xl">
                                <span class="text-transparent bg-clip-text bg-gradient-to-r to-emerald-600 from-sky-400">Protein List</span>
                            </h1>

                            <div class="pb-4 bg-white dark:bg-gray-900">
                                <label for="table-search" class="sr-only">Search</label>
                                <div class="relative mt-1">
                                    <div class="absolute inset-y-0 rtl:inset-r-0 start-0 flex items-center ps-3 pointer-events-none">
                                        <svg class="w-4 h-4 text-gray-500 dark:text-gray-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20">
                                            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"/>
                                        </svg>
                                    </div>
                                    <input wire:model.live="search" type="text" id="table-search" class="block pt-2 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg w-80 bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="Search for imported proteins">
                                </div>
                            </div>

                            <table class="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
                                <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                                    <tr>
                                        <th scope="col" class="p-4">
                                            <div class="flex items-center">
                                                <input id="checkbox-all-search" type="checkbox" class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600">
                                                <label for="checkbox-all-search" class="sr-only">checkbox</label>
                                            </div>
                                        </th>
                                        <th scope="col" class="px-6 py-3">
                                            Protein
                                        </th>
                                        <th scope="col" class="px-6 py-3">
                                            PMIDS
                                        </th>
                                        <th scope="col" class="px-6 py-3">
                                            Category
                                        </th>
                                        <th scope="col" class="px-6 py-3">
                                            Action
                                        </th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr v-for="protein in proteins.data" :key="protein.id" class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600">
                                        <td class="w-4 p-4">
                                            <div class="flex items-center">
                                                <input id="checkbox-table-search-1" type="checkbox" class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600">
                                                <label for="checkbox-table-search-1" class="sr-only">checkbox</label>
                                            </div>
                                        </td>
                                        <th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                                            {{ protein.name }}
                                        </th>
                                        <td class="px-6 py-4">
                                            COUNT
                                        </td>
                                        <td class="px-6 py-4">
                                            Laptop
                                        </td>
                                        <td class="px-6 py-4">
                                            <a class="font-medium text-blue-600 dark:text-blue-500 hover:underline">View</a>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>

                            <div v-if="proteins?.data?.length" class="mt-4">
                                <Pagination :pagination="proteins" />
                            </div>
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </AuthenticatedLayout>
</template>
