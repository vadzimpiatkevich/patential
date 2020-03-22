module ApplicationHelper
  def sidenav_link(path, name)
    html_class = 'active' if current_page?(path)

    content_tag(:li, class: html_class) do
      link_to(name, path)
    end
  end
end
