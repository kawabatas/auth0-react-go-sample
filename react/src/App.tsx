import { useEffect, useState } from 'react'
import { useAuth0 } from "@auth0/auth0-react";
import axios from "axios";
import './App.css'

const useAuth0Token = () => {
  const { isAuthenticated, getAccessTokenSilently } = useAuth0();
  const [accessToken, setAccessToken] = useState("");

  useEffect(() => {
    const fetchToken = async () => {
      setAccessToken(await getAccessTokenSilently())
    };

    if (isAuthenticated) {
      fetchToken();
    }
  }, [isAuthenticated]);

  return accessToken;
};

function App() {
  const [message, setMessage] = useState("");
  const {
    loginWithRedirect,
    isAuthenticated,
    logout,
    user,
  } = useAuth0();
  const token = useAuth0Token();

  const onPublicAPICall = async () => {
    setMessage("");
    const response = await axios({
      url: `${import.meta.env.VITE_API_BASE_URL}/api/public`,
      method: "GET",
      headers: {
        "content-type": "application/json",
      },
    });
    setMessage(JSON.stringify(response.data));
  };

  const onPrivateAPICall = async () => {
    setMessage("");
    const response = await axios({
      url: `${import.meta.env.VITE_API_BASE_URL}/api/private`,
      method: "GET",
      headers: {
        "content-type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    });
    setMessage(JSON.stringify(response.data));
  };

  const onPrivateRBACAPICall = async () => {
    setMessage("");
    const response = await axios({
      url: `${import.meta.env.VITE_API_BASE_URL}/api/privaterbac`,
      method: "GET",
      headers: {
        "content-type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    });
    setMessage(JSON.stringify(response.data));
  };

  return (
    <>
      <h1>Auth0 + Vite + React</h1>
      <div className="card">
        <button onClick={() => loginWithRedirect()} disabled={ isAuthenticated }>
          { isAuthenticated ? "ログイン済み" : "ログイン" }
        </button>
      </div>
      <div className="card">
        {isAuthenticated && (
          <p>{ user?.name }さん、ようこそ</p>
        )}
        {isAuthenticated && (
          <button onClick={() => 
            logout({
              logoutParams: {
                returnTo: window.location.origin,
              }
            })
          }>
            ログアウト
          </button>
        )}
      </div>
      <div className="card">
        <button onClick={() => onPublicAPICall()}>Public API</button>
        <button onClick={() => onPrivateAPICall()}>Private API</button>
        <button onClick={() => onPrivateRBACAPICall()}>Private RBAC API</button>
      </div>
      <div className="card">
        { message }
      </div>
    </>
  )
}

export default App
