const $ = window.jQuery = require('jquery');
const marked = require("marked")
import hljs from 'highlight.js';

marked.setOptions({
  highlight: function (code, lang, _callback) {
    if (hljs.getLanguage(lang)) {
      return hljs.highlight(lang, code).value
    } else {
      return hljs.highlightAuto(code).value
    }
  },
  langPrefix:'hljs p-4 ',
  pedantic: false,
  gfm: true,
  breaks: false,
  sanitize: false,
  smartLists: true,
  smartypants: false,
  xhtml: false
})

window.marked = marked


$(document).ready(function () {
  var preview = $("#preview");
  $("#editor").on("change", function () {
    preview.html(marked($(this).val()));
  }); 

  $(".btn-write").click(function(){
    $(this).addClass("bg-primary text-white");
    $(".btn-preview").removeClass("bg-primary text-white");
    $(".editor-col").show();
    $(".preview").hide();
  }); 

  $(".btn-preview").click(function(){
    $(this).addClass("bg-primary text-white");
    $(".btn-write").removeClass("bg-primary text-white");
    $(".editor-col").hide();
    $(".preview").html(marked($("#editor").val())).show();
  }); 
})