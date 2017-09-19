module Jekyll
  class MarkdownTag < Liquid::Tag
    def initialize(tag_name, text, tokens)
      super
      @text = text.strip
    end
    require "kramdown"
    def render(context)
      "#{Kramdown::Document.new(File.read(File.join(Dir.pwd, '_includes', @text))).to_html}"
    end
  end
end
Liquid::Template.register_tag('markdown', Jekyll::MarkdownTag)
