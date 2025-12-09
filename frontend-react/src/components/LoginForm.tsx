import { useState, type JSX } from "react";
import { type AuthFormProps } from "../App.tsx";

export function AuthenticationForm({ onSubmit }: AuthFormProps): JSX.Element {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [showPassword, setShowPassword] = useState<boolean>(false);
  const [errorMessage, setErrorMessage] = useState<string>("");
  const [showForm, setShowForm] = useState<boolean>(false);

  const togglePasswordVisible = () => {
    setShowPassword(!showPassword);
  };

  const toggleFormShown = () => {
    setShowForm(!showForm);
  };

  return showForm ? (
    <>
      <button onClick={toggleFormShown}>Show login</button>
      <h1>Login</h1>
      <div className="login-fields-div">
        <input
          type="text"
          className="login-fields"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
      </div>
      <div className="login-fields-div">
        <input
          type={showPassword ? "text" : "password"}
          className="login-fields"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </div>
      <div className="buttons-div">
        <button
          onClick={() => {
            onSubmit(username, password, false).then((res) =>
              res.text().then((text) => setErrorMessage(text)),
            );
            setUsername("");
            setPassword("");
          }}
        >
          Login
        </button>
        <button
          className="buttons"
          onClick={() => {
            onSubmit(username, password, true).then((res) =>
              res.text().then((text) => setErrorMessage(text)),
            );
            setUsername("");
            setPassword("");
          }}
        >
          Register
        </button>
        <button
          onClick={() => {
            togglePasswordVisible();
          }}
        >
          {showPassword ? "hide password" : "show password"}
        </button>
      </div>
      <div>{errorMessage}</div>
    </>
  ) : (
    <>
      <button onClick={toggleFormShown}>Show login</button>
    </>
  );
}
