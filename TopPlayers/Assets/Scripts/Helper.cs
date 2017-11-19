using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public static class Helper
{
    public static int GenerateRandomScore()
    {
        return new System.Random().Next(1, 1000);
    }

    public static string ConvertPlayerToString(Player player)
    {
        Dictionary<string, string> pairForJson = new Dictionary<string, string>();
        pairForJson.Add("Name", player.Name);
        pairForJson.Add("Password", player.Password);
        if (player.Score != 0)
        {
            pairForJson.Add("Score", player.Score.ToString());
        }

        JSONObject json = JSONObject.Create(pairForJson);

        return json.ToString();
    }

    public static object GetValueFromJson(JSONObject jSONObject)
    {
        switch (jSONObject.type)
        {
            case JSONObject.Type.STRING:
                {
                    return jSONObject.str;
                }
            case JSONObject.Type.ARRAY:
                {
                    List<Player> players = new List<Player>();

                    foreach (JSONObject item in jSONObject.list)
                    {
                        players.Add(new Player { Name = item["Name"].str, Score = int.Parse(item["Score"].str) });
                    }
                    return players;
                }
        }
        return "";
    }
}
