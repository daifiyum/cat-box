<template>
    <v-snackbar v-model="snackbarVisible" :timeout="2000" color="primary" variant="tonal">
        {{ message }}
        <template v-slot:actions>
            <v-btn color="blue" variant="text" @click="snackbarVisible = false" icon="mdi-close-circle-outline">
            </v-btn>
        </template>
    </v-snackbar>
</template>

<script setup>
import { ref } from "vue"
const snackbarVisible = ref(false)
const message = ref("")

function openSnackbar(msg) {
    message.value = msg;
    snackbarVisible.value = true;
}

defineExpose({
    openSnackbar,
});
</script>

<style lang="css" scoped></style>