ActionView::Base.field_error_proc = Proc.new do |html_tag, instance_tag|
  fragment = Nokogiri::HTML.fragment(html_tag)
  field = fragment.at('input,select,textarea')

  html = if field
           field['class'] = "#{field['class']} invalid"
           html = <<-HTML
              #{fragment.to_s}
              <span class="helper-text">
                #{instance_tag.error_message.to_sentence}
              </span>
           HTML
           html
         else
           html_tag
         end

  html.html_safe
end
