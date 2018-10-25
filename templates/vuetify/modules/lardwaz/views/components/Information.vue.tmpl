<template>
  <v-form>
    <v-container>
      <v-layout row wrap>

        <template v-for="(setting, objectName) in meta">


          <template v-if="useComponent(setting.type, objectName)">
            <v-flex md12>
              <component
                  :is="setting.type"
                  v-model="setting.value"
                  :key="setting.attr.label"
                  :items="setting.attr.items"
                  v-bind="setting.attr"
                  :prepend-icon="setting.attr.icon"
              ></component>
            </v-flex>
          </template>

          <template v-if="setting.type === 'v-radio-group'">
            <v-flex md12>
              <v-radio-group v-model="setting.value" v-bind="setting.attr">
                <v-radio
                    v-for="(radio, index) in setting.values"
                    :key="index"
                    :label="radio.label"
                    :value="radio.value"
                ></v-radio>
              </v-radio-group>
            </v-flex>
          </template>

          <template v-if="setting.type === 'v-date-picker'">
            <v-flex md12>
              <v-menu
                  :ref="setting.ref"
                  :close-on-content-click="false"
                  v-model="setting.menu_value"
                  :nudge-right="40"
                  :return-value.sync="setting.date_value"
                  lazy
                  transition="scale-transition"
                  offset-y
                  full-width
                  min-width="290px"
              >
                <v-text-field
                    slot="activator"
                    v-model="setting.date_value"
                    hint="MM/DD/YYYY format"
                    :label="setting.attr.label"
                    prepend-icon="event"
                    readonly
                ></v-text-field>
                <v-date-picker v-model="setting.date_value" @input="$refs[setting.ref][0].save(setting.date_value)"></v-date-picker>
              </v-menu>
            </v-flex>
          </template>

          <template v-if="setting.type === 'v-combobox'">
            <v-flex md12>
              <v-combobox
                  v-model="setting.value"
                  :items="setting.attr.items"
                  :label="setting.attr.label"
                  multiple
                  chips
                  :prepend-icon="setting.attr.icon"
              ></v-combobox>
            </v-flex>
          </template>

          <template v-if="setting.type === 'misc'">
            <v-flex md12>
              <v-card>
                <v-layout row wrap pa-4>
                  <template v-for="(item, objectName) in setting.options">
                    <v-flex md6>
                      <component
                          :is="item.type"
                          v-model="item.value"
                          :key="item.attr.label"
                          v-bind="item.attr"
                      ></component>
                    </v-flex>
                  </template>
                </v-layout>
              </v-card>
            </v-flex>

          </template>
        </template>
      </v-layout>
    </v-container>
  </v-form>
</template>
<script>
  import { mapGetters } from "vuex";

 export default {
   data () {
     return {
       modal: false
     }
   },
   computed: {
     ...mapGetters(['meta'])
   },
   methods: {
     useComponent(component, objectName) {

       console.log(objectName)
       if (
           component === 'v-radio-group' ||
           component === 'v-date-picker' ||
           component === 'v-combobox' ||
           component === 'misc'
       ) {
         return false
       } else {
         return true
       }
     }
   },
   component:{

   },
   mounted () {
     console.log(this.$refs)
   }
 }
</script>
<style lang="scss" scoped>

</style>