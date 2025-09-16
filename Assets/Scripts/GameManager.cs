using UnityEngine;
using System.Collections;
using System.Collections.Generic; // Added for List
using System.Linq; // Added for Linq operations like ToList()

public class GameManager : MonoBehaviour
{
    public Sprite monsterSprite; // Public variable to assign monster sprite in editor

    // Player Resources
    public int playerGold = 100;
    public int playerLives = 10;

    // Wave Management
    public int waveNumber = 0;
    public int monstersPerWave = 5;
    public float genInterval = 1f;
    public float timeBetweenWaves = 5f;
    public int currentWaveExperience = 50; // Experience awarded per wave

    private int monstersGenedThisWave = 0;
    private bool waveInProgress = false;
    private List<Monster> activeMonsters = new List<Monster>(); // Track active monsters

    private MapManager _mapManager;

    // Monster Prefab (will be set up in editor)
    public GameObject monsterPrefab; // Assign the Monster GameObject as a prefab

    void Awake()
    {
        _mapManager = FindAnyObjectByType<MapManager>();
        if (_mapManager == null)
        {
            Debug.LogError("MapManager not found in the scene! Please add a MapManager GameObject.");
        }
    }

    void OnEnable()
    {
        Monster.OnMonsterDied += HandleMonsterDied;
    }

    void OnDisable()
    {
        Monster.OnMonsterDied -= HandleMonsterDied;
    }

    void Start()
    {
        // Map loading will now be triggered by StartGameWithMap after user selection
        // For testing, you might temporarily call StartGameWithMap(0) here
        // StartGameWithMap(0); // Example: automatically start with the first map for testing

        // Create a simple square sprite programmatically for Hero
        Texture2D texture = new Texture2D(1, 1);
        texture.SetPixel(0, 0, Color.white);
        texture.Apply();
        Sprite defaultSquareSprite = Sprite.Create(texture, new Rect(0, 0, texture.width, texture.height), Vector2.one * 0.5f);

        // Create Hero GameObject
        GameObject hero = new GameObject("Hero");
        hero.transform.position = new Vector3(0f, 0f, 0f); // X, Y, Z (Z=0 for 2D)
        hero.transform.localScale = new Vector3(1f, 1f, 1f); // Set a visible size
        SpriteRenderer heroRenderer = hero.AddComponent<SpriteRenderer>();
        heroRenderer.sprite = defaultSquareSprite; // Assign programmatically created sprite
        heroRenderer.color = Color.blue; // Make hero blue
        hero.AddComponent<Hero>(); // Add the Hero script component
        hero.AddComponent<BoxCollider2D>(); // Add a BoxCollider2D for potential future physics interactions

        // Instead of creating monster directly, we'll use a prefab for gening
        // For now, let's create a temporary monster GameObject to make it a prefab
        GameObject tempMonster = new GameObject("Monster");
        // tempMonster.transform.position = _mapManager.MonsterStartPoint; // This will be set when map is loaded
        tempMonster.transform.localScale = new Vector3(32f, 32f, 1f);
        SpriteRenderer tempMonsterRenderer = tempMonster.AddComponent<SpriteRenderer>();
        tempMonsterRenderer.sprite = monsterSprite;
        tempMonster.AddComponent<BoxCollider2D>();
        tempMonster.AddComponent<MonsterMovement>(); // Add MonsterMovement to the temp monster
        tempMonster.AddComponent<Monster>(); // Add Monster to the temp monster
        // We will make this a prefab in the editor and assign it to monsterPrefab
        // Then destroy this temp monster
        Destroy(tempMonster);

        // Start waves only after a map has been selected and loaded
        // StartCoroutine(StartWaves());
    }

    // How to integrate UI for map selection and game start:
    // 1. Create UI elements (e.g., Buttons, Dropdown) in a Canvas in your scene.
    // 2. For each UI element that triggers game start with a specific map:
    //    - Select the UI element in the Hierarchy.
    //    - In the Inspector, find its event (e.g., Button's OnClick(), Dropdown's OnValueChanged).
    //    - Click the '+' button to add a new event.
    //    - Drag the GameObject with this GameManager script attached into the event slot.
    //    - From the dropdown menu, select `GameManager` -> `StartGameWithMap(int mapIndex)` or `GameManager` -> `StartGameWithMap(string mapName)`.
    //    - Enter the desired `mapIndex` (based on the order in MapManager's `availableMaps` list) or `mapName`.
    // This allows the player to select a map and then start the game with that selection.
    public void StartGameWithMap(int mapIndex)
    {
        if (_mapManager == null)
        {
            Debug.LogError("MapManager not found. Cannot start game.");
            return;
        }

        _mapManager.SelectAndLoadMap(mapIndex);

        // Now that the map is loaded, we can start the waves
        StartCoroutine(StartWaves());
    }

    public void StartGameWithMap(string mapName)
    {
        if (_mapManager == null)
        {
            Debug.LogError("MapManager not found. Cannot start game.");
            return;
        }

        _mapManager.SelectAndLoadMap(mapName);

        // Now that the map is loaded, we can start the waves
        StartCoroutine(StartWaves());
    }

    IEnumerator StartWaves()
    {
        while (playerLives > 0) // Game continues as long as player has lives
        {
            yield return new WaitForSeconds(timeBetweenWaves);
            waveNumber++;
            Debug.Log("Starting Wave " + waveNumber);
            monstersGenedThisWave = 0;
            waveInProgress = true;
            activeMonsters.Clear(); // Clear monsters from previous wave
            StartCoroutine(GenMonstersInWave());

            // Wait until all monsters for this wave are gened AND all active monsters are defeated
            yield return new WaitWhile(() => waveInProgress || activeMonsters.Count > 0);

            // Wave completed, award experience
            AwardExperienceToHeroes();
        }
        Debug.Log("Game Over!");
    }

    IEnumerator GenMonstersInWave()
    {
        for (int i = 0; i < monstersPerWave; i++)
        {
            GenMonster();
            monstersGenedThisWave++;
            yield return new WaitForSeconds(genInterval);
        }
        waveInProgress = false; // All monsters for this wave have been gened
    }

    void GenMonster()
    {
        if (monsterPrefab == null)
        {
            Debug.LogError("Monster Prefab is not assigned in GameManager!");
            return;
        }

        // Instantiate monster from prefab
        GameObject newMonsterGO = Instantiate(monsterPrefab, _mapManager.MonsterStartPoint, Quaternion.identity);
        Monster newMonster = newMonsterGO.GetComponent<Monster>();
        if (newMonster != null)
        {
            activeMonsters.Add(newMonster); // Add to active monsters list
        }

        MonsterMovement monsterMovement = newMonsterGO.GetComponent<MonsterMovement>();
        if (monsterMovement != null)
        {
            monsterMovement.start = _mapManager.MonsterStartPoint;
            monsterMovement.end = _mapManager.MonsterEndPoint;
        }
    }

    private void HandleMonsterDied(Monster monster)
    {
        if (activeMonsters.Contains(monster))
        {
            activeMonsters.Remove(monster);
            Debug.Log($"Monster {monster.name} removed. {activeMonsters.Count} monsters remaining in wave.");
        }
    }

    private void AwardExperienceToHeroes()
    {
        Hero[] activeHeroes = FindObjectsByType<Hero>(FindObjectsSortMode.None);
        foreach (Hero hero in activeHeroes)
        {
            hero.GainExperience(currentWaveExperience);
        }
        Debug.Log($"Awarded {currentWaveExperience} experience to all active heroes.");
    }

    // Update is called once per frame
    void Update()
    {

    }
}

