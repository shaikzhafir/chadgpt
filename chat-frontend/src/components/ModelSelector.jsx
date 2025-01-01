/* eslint-disable react/prop-types */
export const ModelSelector = ({ selectedModel, onModelChange, disabled = false }) => {
    const models = [
        { value: 'openai', label: 'GPT-4 (OpenAI)' },
        { value: 'anthropic', label: 'Claude something' },
        { value: 'deepseek', label: 'DeepSeek' },
    ];

    return (
        <select
            value={selectedModel}
            onChange={(e) => onModelChange(e.target.value)}
            disabled={disabled}
            className="w-full p-2 border rounded-lg bg-white focus:ring-2 focus:ring-blue-500"
        >
            {models.map((model) => (
                <option key={model.value} value={model.value}>
                    {model.label}
                </option>
            ))}
        </select>
    );
};


export default ModelSelector;