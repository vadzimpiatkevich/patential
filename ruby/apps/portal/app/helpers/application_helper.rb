module ApplicationHelper
  def sidenav_link(path, name, method: :get)
    html_class = 'active' if current_page?(path)

    content_tag(:li, class: html_class) do
      link_to(name, path, method: method)
    end
  end
end
