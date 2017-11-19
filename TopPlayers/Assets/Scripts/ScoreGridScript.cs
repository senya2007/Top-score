using System.Collections;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;
using UnityEngine.UI;

public class ScoreGridScript : MonoBehaviour {

    public GameObject playerScoreEntryPrefab;

    ScoreManager scoreManager;

    int lastChangeCounter;

    // Use this for initialization
    void Start()
    {
        scoreManager = GameObject.FindObjectOfType<ScoreManager>();

        lastChangeCounter = scoreManager.GetChangeCounter();
    }

    // Update is called once per frame
    void Update()
    {
        if (scoreManager == null)
        {
            Debug.LogError("You forgot to add the score manager component to a game object!");
            return;
        }

        if (scoreManager.GetChangeCounter() == lastChangeCounter)
        {
            // No change since last update!
            return;
        }

        lastChangeCounter = scoreManager.GetChangeCounter();

        while (this.transform.childCount > 0)
        {
            Transform c = this.transform.GetChild(0);
            c.SetParent(null);  // Become Batman
            Destroy(c.gameObject);
        }

        

        foreach (var player in scoreManager.GetAllPlayers())
        {
            GameObject go = (GameObject)Instantiate(playerScoreEntryPrefab);
            go.transform.SetParent(this.transform);
            go.transform.localScale = new Vector3(1, 1, 1);
            go.transform.Find("Name").GetComponent<Text>().text = player.Key;
            go.transform.Find("Score").GetComponent<Text>().text = scoreManager.GetScore(player.Key).ToString();
        }
    }

}
