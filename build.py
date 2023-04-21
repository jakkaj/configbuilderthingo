import yaml
from jinja2 import Template

class ToolParser:
    def __init__(self, yaml_file):
        self.yaml_data = self.load_yaml_file(yaml_file)

    def load_yaml_file(self, file_path):
        try:
            with open(file_path, 'r') as yaml_file:
                yaml_content = yaml_file.read()
                return yaml.safe_load(yaml_content)
        except Exception as e:
            print("Error reading YAML file:", e)
            return None

    def get_tool_by_type(self, tool_type):
        if self.yaml_data:
            for tool in self.yaml_data:
                if tool['type'] == tool_type:
                    return tool
        return None

    def render_template(self, tool, template_args):
        action_section = tool.get('action', {})
        args_template = action_section.get('args_template', '')
        template = Template(args_template)
        rendered_args = template.render(**template_args)

        if action_section.get('type') == 'exec':
            path = action_section.get('path')
            return f"{path} {rendered_args}"
        else:
            return rendered_args

class ArtefactParser:
    def __init__(self, yaml_file):
        self.yaml_data = self.load_yaml_file(yaml_file)

    def load_yaml_file(self, file_path):
        try:
            with open(file_path, 'r') as yaml_file:
                yaml_content = yaml_file.read()
                return yaml.safe_load(yaml_content)
        except Exception as e:
            print("Error reading YAML file:", e)
            return None

    def get_all_tool_params(self):
        tools = {}
        if self.yaml_data:
            tool_list = self.yaml_data.get('tools', [])
            for tool in tool_list:
                tool_type = tool['type']
                params = {}
                for param in tool.get('params', []):
                    params[param['name']] = param['value']
                tools[tool_type] = params
        return tools

# Usage example
artefact_yaml_file = 'artefact_config.yaml'
artefact_parser = ArtefactParser(artefact_yaml_file)
all_tool_params = artefact_parser.get_all_tool_params()

tool_config_yaml_file = 'tool_config.yaml'
tool_parser = ToolParser(tool_config_yaml_file)

bash_commands = []

for tool_type, params in all_tool_params.items():
    tool = tool_parser.get_tool_by_type(tool_type)
    if tool:
        rendered_template = tool_parser.render_template(tool, params)
        print(f"Rendered Template for {tool_type}: {rendered_template}")
        bash_commands.append(rendered_template)
        
print("\nBash Commands:")
for cmd in bash_commands:
    print(cmd)
