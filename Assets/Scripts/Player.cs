using UnityEngine;
using System.Collections.Generic;

public class Player : MonoBehaviour
{
    [System.Serializable]
    private class HeroDataWrapper
    {
        public string[] heroNames;
    }

    private List<MapData> Heros = new List<MapData>(); // Initialize the list

    void Awake()
    {
        LoadHeroData("hero_data"); // Assuming a file named hero_data.json (without extension) in Resources
    }

    private void LoadHeroData(string fileNameWithoutExtension)
    {
        /*
         * To make this code work in your Unity project, please follow these steps:
         *
         * 1. Create a `Resources` Folder: If you don't already have one, create a folder named `Resources`
         *    directly under your `Assets` folder in your Unity project. (e.g., `Assets/Resources`).
         *
         * 2. Create `MapData` ScriptableObjects:
         *    - In the Unity Editor, go to `Assets` -> `Create` -> `Game` -> `Map Data`.
         *    - Create several `MapData` assets (e.g., "Hero1Map", "Hero2Map").
         *    - Save these `MapData` assets inside the `Resources` folder (or any subfolder within `Resources`).
         *
         * 3. Create `hero_data.json`:
         *    - Create a new text file named `hero_data.json` inside your `Resources` folder
         *      (e.g., `Assets/Resources/hero_data.json`).
         *
         * 4. Populate `hero_data.json`:
         *    - Open `hero_data.json` and enter the *names* of the `MapData` assets you created, in a JSON array.
         *      For example:
         *      ```json
         *      {
         *          "heroNames": [
         *              "Hero1Map",
         *              "Hero2Map",
         *              "AnotherHeroMap"
         *          ]
         *      }
         *      ```
         */
        TextAsset heroDataJsonFile = Resources.Load<TextAsset>(fileNameWithoutExtension);

        if (heroDataJsonFile != null)
        {
            HeroDataWrapper dataWrapper = JsonUtility.FromJson<HeroDataWrapper>(heroDataJsonFile.text);

            if (dataWrapper != null && dataWrapper.heroNames != null)
            {
                foreach (string heroName in dataWrapper.heroNames)
                {
                    // Load MapData ScriptableObject by name from Resources
                    MapData heroMapData = Resources.Load<MapData>(heroName.Trim());
                    if (heroMapData != null)
                    {
                        Heros.Add(heroMapData);
                        Debug.Log($"Loaded hero: {heroName.Trim()}");
                    }
                    else
                    {
                        Debug.LogWarning($"Could not load MapData asset: {heroName.Trim()}. Make sure it's in a Resources folder.");
                    }
                }
            }
            else
            {
                Debug.LogError($"Failed to parse hero data from JSON or heroNames array is null/empty in: {fileNameWithoutExtension}");
            }
        }
        else
        {
            Debug.LogError($"Hero data JSON file not found in Resources: {fileNameWithoutExtension}.json");
        }
    }
}
