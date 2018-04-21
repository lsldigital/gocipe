<template>
    <v-navigation-drawer v-bind:value="drawer" clipped fixed app>
        <v-list dense>
            <v-list-tile :to="{ name: 'gocipe'}">
                <v-list-tile-action>
                    <v-icon>build</v-icon>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Builder</v-list-tile-title>
                </v-list-tile-content>
            </v-list-tile>
        </v-list>

        <v-divider></v-divider>
    </v-navigation-drawer>
</template>
