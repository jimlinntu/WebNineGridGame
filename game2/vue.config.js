module.exports = {
    chainWebpack: config => {
        config.module.rules.delete('eslint');
        // https://stackoverflow.com/questions/54951082/vue-cli-title-option-for-htmlwebpackplugin-does-not-work
        config.plugin('html').tap((args)=>{
            args[0].title = "生命助道會退修會"
            return args;
        });
    }
}
