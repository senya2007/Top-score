using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;
using UnityEngine.UI;

public class MainManagerScript : MonoBehaviour {

    public Text loginText;
    public InputField passwordText;
    public Text mainPanelErrorText;
    public Text tableScoreErrorText;

    public GameObject mainPanel;
    public GameObject scorePanel;

    public Button updateScoreButton;
    public ScoreManager scoreManager;

    public bool isConnectedEthernet { get { return Application.internetReachability != NetworkReachability.NotReachable; } }

    private const string SERVER = "http://localhost:8080/";
    private const string LOGIN = "login";
    private const string CREATEPLAYER = "create";
    private const string UPDATESCORE = "updateScore";
    private const string UPDATEALLSCORES = "updateAllScores";

    private const int MAXPLAYERSINTABLE = 6;

    private static Player CurrentPlayer;
    public void Login()
    {
        SendPlayerToServer(LOGIN);
    }

    public void Register()
    {
        SendPlayerToServer(CREATEPLAYER);
    }

    public void SendPlayerToServer(string addres)
    {
        if (CheckForEmptyTextView())
        {
            return;
        }
        CurrentPlayer = new Player { Name = loginText.text, Password = passwordText.text };
        string stringJson = Helper.ConvertPlayerToString(CurrentPlayer);
        StartCoroutine(SendJson(addres, stringJson));
    }

   

    public void UpdateScore()
    {
        if (CurrentPlayer != null)
        {
            if (updateScoreButton.IsInteractable())
            {
                updateScoreButton.interactable = false;
            }
            CurrentPlayer.Score = Helper.GenerateRandomScore();
            string stringJson = Helper.ConvertPlayerToString(CurrentPlayer);
            StartCoroutine(SendJson(UPDATESCORE, stringJson));
        }
    }

    public void UpdateAllScores()
    {
        StartCoroutine(SendJson(UPDATEALLSCORES));
    }

    private IEnumerator SendJson(string methodForSend, string stringJson = "empty")
    {
        if (CheckInternet())
        {
            yield return null;
        }

        WWW w = new WWW(SERVER + methodForSend, System.Text.Encoding.UTF8.GetBytes(stringJson.ToCharArray()));
        yield return w;

        if (w.isDone)
        {
            DeactivateNewScoreButton();

            if (!CheckServerConnect(w.error))
            {
                yield return null;
            }
            else
            {
                ParseJsonResponse(w.text);
                Debug.Log("Json send");
            }
        }
    }

    private bool CheckServerConnect(string error)
    {
        if (!string.IsNullOrEmpty(error) && error == "Cannot connect to destination host")
        {
            mainPanelErrorText.text = "Сервер недоступен";
            tableScoreErrorText.text = "Сервер недоступен";
            if (updateScoreButton.IsInteractable())
            {
                updateScoreButton.interactable = false;
            }
            return false;
        }
        return true;
    }

    private void DeactivateNewScoreButton()
    {
        if (updateScoreButton.IsActive() && !updateScoreButton.IsInteractable())
        {
            updateScoreButton.interactable = true;
        }
    }

    private bool CheckInternet()
    {
        if (!isConnectedEthernet)
        {
            mainPanelErrorText.text = "Нет интернета";
            tableScoreErrorText.text = "Нет интернета";
            return true;
        }
        else
        {
            mainPanelErrorText.text = "";
            tableScoreErrorText.text = "";
        }
        return false;
    }

    private void ParseJsonResponse(string responseJson)
    {
        if (string.IsNullOrEmpty(responseJson))
        {
            mainPanelErrorText.text = "Пустой ответ от сервера";
            return;
        }
        else
        {
            var json = new JSONObject(responseJson);
            var type = json["Type"].str;
            var value = Helper.GetValueFromJson(json["Value"]);
            var methodName = json["MethodName"].str;
            SetActionFromJson(type, value, methodName);
        }
    }

    private IEnumerator UpdateAllScoresCoroutine()
    {
        while (true)
        {
            UpdateAllScores();
            yield return new WaitForSeconds(30);
        }
    }
    private void SetActionFromJson(string type, object value, string methodName)
    {
        switch (methodName)
        {
            case "login":
            case "create":
                {
                    if (type == "good")
                    {
                        mainPanel.SetActive(false);
                        scorePanel.SetActive(true);
                        StartCoroutine(UpdateAllScoresCoroutine());
                    }
                    else if (type == "error")
                    {
                        CurrentPlayer = null;
                        mainPanelErrorText.text = (string)value;
                    }
                    else
                    {
                        mainPanelErrorText.text = "Неизвестная ошибка";
                    }
                }
                break;
            case "updateScore":
                if (type == "good")
                {
                    StopAllCoroutines();
                    StartCoroutine(UpdateAllScoresCoroutine());
                }
                else if (type == "error")
                {
                    CurrentPlayer = null;
                    tableScoreErrorText.text = (string)value;
                }
                else
                {
                    tableScoreErrorText.text = "Неизвестная ошибка";
                }
                break;
            case "updateAllScores":
                {
                    if (type == "players")
                    {
                        scoreManager.Reset();
                        int i = 1;
                        foreach (var player in (value as List<Player>).OrderByDescending(x => x.Score).ToList())
                        {
                            if (i > MAXPLAYERSINTABLE)
                            {
                                continue;
                            }
                            scoreManager.SetScore(player.Name, player.Score);
                            i++;
                        }
                    }
                    else if (type == "error")
                    {
                        tableScoreErrorText.text = (string)value;
                    }
                    else
                    {
                        tableScoreErrorText.text = "Неизвестная ошибка";
                    }
                }
                break;
        }
    }

    private bool CheckForEmptyTextView()
    {
        if (string.IsNullOrEmpty(loginText.text) && string.IsNullOrEmpty(passwordText.text))
        {
            mainPanelErrorText.text = "Введите логин и пароль";
            return true;
        }
        if (string.IsNullOrEmpty(loginText.text))
        {
            mainPanelErrorText.text = "Введите логин";
            return true;
        }

        if (string.IsNullOrEmpty(passwordText.text))
        {
            mainPanelErrorText.text = "Введите пароль";
            return true;
        }

        return false;
    }
}
