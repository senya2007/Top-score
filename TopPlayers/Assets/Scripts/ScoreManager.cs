using UnityEngine;
using System.Collections.Generic;
using System.Linq;

public class ScoreManager : MonoBehaviour
{

    Dictionary<string, int> playerScores;

    int changeCounter = 0;

    void Start()
    {
    }

    void Init()
    {
        if (playerScores != null)
            return;

        playerScores = new Dictionary<string, int>();
    }

    public void Reset()
    {
        changeCounter++;
        playerScores = null;
    }

    public int GetScore(string username)
    {
        Init();

        if (playerScores.ContainsKey(username) == false)
        {
            // We have no score record at all for this username
            return 0;
        }

        return playerScores[username];
    }

    public void SetScore(string username, int value)
    {
        Init();

        changeCounter++;

        if (playerScores.ContainsKey(username) == false)
        {
            playerScores[username] = 0;
        }

        playerScores[username] = value;
    }

    public void ChangeScore(string username, int amount)
    {
        Init();
        int currScore = GetScore(username);
        SetScore(username, currScore + amount);
    }

    public string[] GetPlayerNames()
    {
        Init();
        return playerScores.Keys.ToArray();
    }


    public Dictionary<string, int> GetAllPlayers()
    {
        return playerScores;
    }


    public int GetChangeCounter()
    {
        return changeCounter;
    }

    public void DEBUG_ADD_KILL_TO_QUILL()
    {
        ChangeScore("quill18", 1);
    }

    public void DEBUG_INITIAL_SETUP()
    {
        SetScore("quill18", 0);
        SetScore("quill18", 345);

        SetScore("bob", 1000);
        SetScore("bob", 14345);

        SetScore("AAAAAA", 3);
        SetScore("BBBBBB", 2);
        SetScore("CCCCCC", 1);


        Debug.Log(GetScore("quill18"));
    }

}
