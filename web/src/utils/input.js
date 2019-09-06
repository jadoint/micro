export const handleTextChange = ({ event, inputState, setInputState }) => {
  const { name, value } = event.target;

  setInputState({ ...inputState, [name]: value });
};

export const handleCheckboxChange = ({ event, inputState, setInputState }) => {
  const { name, checked } = event.target;

  setInputState({ ...inputState, [name]: checked });
};
