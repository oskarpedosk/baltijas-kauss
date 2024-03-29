import fs from 'fs'
import puppeteer from 'puppeteer-extra'
import StealthPlugin from 'puppeteer-extra-plugin-stealth'
import AdblockerPlugin from 'puppeteer-extra-plugin-adblocker';
const systemOS = process.argv[2]
const playerID = process.argv[3]
const ratingsURL = process.argv[4]
let page
let imagesPath
let browserPath


if (systemOS === "mac") {
    browserPath = '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'
    imagesPath = '../../images/badges'
} else if (systemOS === "windows") {
    browserPath = 'C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe';
    imagesPath = '.\\static\\images\\badges';    
} else if (systemOS === "ubuntu") {
    browserPath = '/usr/bin/chromium-browser'
    imagesPath = '/var/www/bkauss/static/images/badges'
}

puppeteer.use(StealthPlugin())
puppeteer.use(AdblockerPlugin())
puppeteer.launch({
    executablePath: browserPath, 
    headless: "new", 
    args: ['--no-sandbox', '--disable-setuid-sandbox'],
    timeout: 40000
    }).then(async browser => {
    page = await browser.newPage()

    const player_and_badges = await scrapePlayer(ratingsURL);
    console.log(JSON.stringify(player_and_badges));
    await browser.close();
})


async function scrapePlayer(url) {
    await page.goto(url) 

    const player = {
        player_id: parseInt(playerID),
        first_name: null,
        last_name: null,
        primary_position: null,
        secondary_position: null,
        team_id: null,
        assigned_position: null,
        archetype: null,
        height: null,
        weight: null,
        nba_team: null,
        nationality: null,
        birthdate: null,
        jersey: null,
        draft: null,
        img_url: null,
        ratings_url: url,
        overall: 0,
        attributes: null,
        bronze_badges: 0,
        silver_badges: 0,
        gold_badges: 0,
        hof_badges: 0,
        total_badges: 0,
    }

    // Name
    const name = await page.evaluate(() =>
        document.querySelector('.header-title').innerText);
    let regex = /(\S+)\s/;
    if (regex.test(name)) {
        player.first_name = getGroup(regex, name, 1)
        regex = /\S+\s(.+\S|)/;
        player.last_name = getGroup(regex, name, 1)
    }

    // Info (team, height, weight etc.)
    const player_info = await page.evaluate(() =>
            (Array.from(document.querySelectorAll('.header-subtitle > p'))
            .map(element => element.innerText)));

    for (const element of player_info) {
        let info = element;
        
        // NBA team
        let regex = /Team: (.+\w)/;
        if (regex.test(info)) {
            player.nba_team = getGroup(regex, info, 1);
            continue;
        }
        // Archetype
        regex = /Archetype: (.+\w)/;
        if (regex.test(info)) {
            player.archetype = getGroup(regex, info, 1);
            continue;
        }
        // Primary position
        regex = /Position: (\w+)/;
        if (regex.test(info)) {
            player.primary_position = getGroup(regex, info, 1);
            // Secondary position
            regex = /\/ (\w+)/;
            if (regex.test(info)) {
                player.secondary_position = getGroup(regex, info, 1);
            }
            continue;
        }
        // Height
        regex = /\((\d+)cm\)/;
        if (regex.test(info)) {
            player.height = parseInt(getGroup(regex, info, 1));
            regex = /\((\d+)kg\)/;
            if (regex.test(info)) {
                player.weight = parseInt(getGroup(regex, info, 1));
            }
            continue;
        }
        // Jersey
        regex = /Jersey: (#\d+)$/;
        if (regex.test(info)) {
            player.jersey = getGroup(regex, info, 1);
        }
    }

    // Add overall rating
    const overall_rating = await page.evaluate(() => {
        const overall_rating_element = document.querySelector('.attribute-box-player');
        return overall_rating_element ? parseInt(overall_rating_element.innerText) : null;
        });
        
    if (overall_rating) {
    player.overall = overall_rating;
    }

    // Scrape attributes
    const attributes_scrape = await page.evaluate(() => {
        const attributes_element = document.getElementById('nav-attributes');
        if (!attributes_element) {
          return [];
        }
        return [
          ...attributes_element.querySelectorAll('.card-header'),
          ...attributes_element.querySelectorAll('li')
        ].map(element => element.innerText);
      });
      
    // Add attributes
    const attributes = addAttributes(attributes_scrape)
    player.attributes = attributes;

    // Scrape player badges
    const badges = await scrapePlayerBadges();

    for (const badge of badges) {
        const level = getGroup(/[-_](\w+)\./, badge.url, 1)
        if (level === "bronze" || level === "Bronze") {
            player.bronze_badges++
        } else if (level === "silver" || level === "Silver") {
            player.silver_badges++
        } else if (level === "gold" || level === "Gold") {
            player.gold_badges++
        } else if (level === "hof" || level === "HOF") {
            player.hof_badges++
        }
        player.total_badges++
    }

    // Scrape nba.com info and image
    const nba_com_scrape = await scrapeNBAcom(player.first_name, player.last_name);

    Object.assign(player, nba_com_scrape);

    return [player, badges]
}

async function scrapeNBAcom(first_name, last_name) {
    const nba_com_scrape = {};
    try {
        // Search for player with DuckDuckGo
        const query = `site:nba.com/player ${first_name}+${last_name}`;
        await page.goto('https://www.duckduckgo.com/?q=' + query);
        await page.waitForSelector('#r1-0 > div:nth-child(2) > h2:nth-child(1) > a:nth-child(1)');
        const player_url = await page.evaluate(() => document.querySelector('#r1-0 > div:nth-child(2) > h2:nth-child(1) > a:nth-child(1)').href);

        // Go to nba.com
        await page.goto(player_url);

        // Get player info
        const nba_com_info = await page.evaluate(() => {
            const parent_elements = document.querySelectorAll('div.PlayerSummary_hw__HNuGb > div > div');
            if (parent_elements) {
                return Array.from(parent_elements).map(parent => 
                    Array.from(parent.querySelectorAll('p')).map(element => element.innerText)
                );
            } else {
                return null;
            }
        });
        for (const element of nba_com_info) {
            if (element.length === 2) {
                const info = element[0]
                const data = element[1]
                if (info === "HEIGHT") {
                    let regex = /\((\d\d\d)m\)/;
                    const height = data.replace('.', '');
                    if (regex.test(height)) {
                        nba_com_scrape.height = parseInt(getGroup(regex, height, 1));
                    }
                } else if (info === "WEIGHT") {
                    let regex = /\((\d+)kg\)/;
                    if (regex.test(data)) {
                        nba_com_scrape.weight = parseInt(getGroup(regex, data, 1));
                    }
                } else if (info === "COUNTRY") {
                    nba_com_scrape.nationality = data;
                } else if (info === "DRAFT") {
                    nba_com_scrape.draft = data;
                } else if (info === "BIRTHDATE") {
                    nba_com_scrape.birthdate = data;
                }

            }
        }

        // Scrape jersey number
        if (await page.$('p.PlayerSummary_mainInnerInfo__jv3LO') !== null) {
            const nba_dot_com_jersey = await page.evaluate(() => document.querySelector('p.PlayerSummary_mainInnerInfo__jv3LO').innerText);
            let regex = /(#\d+)/;
            if (regex.test(nba_dot_com_jersey)) {
                nba_com_scrape.jersey = getGroup(regex, nba_dot_com_jersey, 1);

            }
        }

        // Scrape player image url
        if (await page.$('img.PlayerImage_image__wH_YX:nth-child(2)') !== null) {
            nba_com_scrape.img_url = await page.evaluate(() => document.querySelector('img.PlayerImage_image__wH_YX:nth-child(2)').src);
        } else if (await page.$('img.PlayerImage_image__wH_YX') !== null) {
                nba_com_scrape.img_url = await page.evaluate(() => document.querySelector('img.PlayerImage_image__wH_YX').src);
            } else {
                nba_com_scrape.img_url = null;
            }

        return nba_com_scrape
    } catch {
        console.log("Couldn't get player data from nba.com: " + first_name + " " + last_name);
        return nba_com_scrape
    }
}

function addAttributes(attributes_scrape) {
    let attributes = {
        OutsideScoring: null,
        Athleticism: null,
        InsideScoring: null,
        Playmaking: null,
        Defending: null,
        Rebounding: null,
        Intangibles: null,
        Potential: null,
        TotalAttributes: null,
        CloseShot: null,
        MidRangeShot: null,
        ThreePointShot: null,
        FreeThrow: null,
        ShotIQ: null,
        OffensiveConsistency: null,
        Speed: null,
        Acceleration: null,
        Strength: null,
        Vertical: null,
        Stamina: null,
        Hustle: null,
        OverallDurability: null,
        Layup: null,
        StandingDunk: null,
        DrivingDunk: null,
        PostHook: null,
        PostFade: null,
        PostControl: null,
        DrawFoul: null,
        Hands: null,
        PassAccuracy: null,
        BallHandle: null,
        SpeedwithBall: null,
        PassIQ: null,
        PassVision: null,
        InteriorDefense: null,
        PerimeterDefense: null,
        Steal: null,
        Block: null,
        LateralQuickness: null,
        HelpDefenseIQ: null,
        PassPerception: null,
        DefensiveConsistency: null,
        OffensiveRebound: null,
        DefensiveRebound: null,
    }
    
    for (const element of attributes_scrape) {
        let regex = /(\S+) (.+\w)/;
        let attribute_str = getGroup(regex, element, 2);
        attribute_str = attribute_str.replaceAll(' ', '');
        attribute_str = attribute_str.replaceAll('-', '');
        let attribute_int = null;
        if (attribute_str === "TotalAttributes") {
            regex = /\d+/g;
            attribute_int = parseInt(getGroup(regex, element, 0) + getGroup(regex, element, 1));
            if (attribute_int === 0) {
                attribute_int = null;
            }
        } else {
            attribute_int = parseInt(getGroup(regex, element, 1));
        }
        attributes[attribute_str] = attribute_int;
    }

    // Replace NaN attributes with null
    for (let key in attributes) {
        if (isNaN(attributes[key])) {
            attributes[key] = null;
        }
    }
    return attributes
}

function getGroup(regex, str, index) {
    try {
        return str.match(regex)[index]
    } catch {
        return null
    }
}

async function scrapePlayerBadges() {
    const badges_scrape = await page.evaluate(() => {
        const badges = document.getElementById('pills-all');
        if (!badges) {
            return []
        }
        return Array.from(badges.querySelectorAll('.card-body')).map(element => 
            Array.from(element.querySelectorAll('*')).map(el => el.innerText))
    });

    const badges_levels = await page.evaluate(() => {
        const badges = document.getElementById('pills-all');
        if (!badges) {
            return []
        } 
        return Array.from(badges.querySelectorAll('img')).map(element => element.getAttribute('data-src'))
    });

    const badges = []

    for (let i = 0; i < badges_scrape.length; i++) {
        const name = badges_scrape[i][0]
        const type = badges_scrape[i][1]
        const info = badges_scrape[i][2]
        const badge = {
            name: name,
            type: type,
            info: info,
            img_id: await saveImage(badges_levels[i], imagesPath, getGroup(/uploads\/(.+)\.png/, badges_levels[i], 1)),
            level: getGroup(/_(\w+)\./, badges_levels[i], 1),
            url: badges_levels[i],
        }
        badges.push(badge)
    }
    return badges
}

async function saveImage(url, folder, file_name) {
    try {
      let extension = '.png';
      let regex = /(\.png|\.webp)$/;
      if (regex.test(url)) {
        extension = getGroup(regex, url, 1);
      }
      const filePath = folder + '/' + (file_name + extension);
      const response = await page.goto(url);
      const buffer = await response.buffer();
  
      // Check if file exists
      if (fs.existsSync(filePath)) {
        return `${file_name}${extension}`;
      }
  
      fs.writeFileSync(filePath, buffer);
      return `${file_name}${extension}`;
    } catch (error) {
      console.error(`${folder} ${file_name} couldn't save image: ${error}`);
      return 'default.png';
    }
  }