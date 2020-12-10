// ExpressJS Setup
const express = require('express');
const app = express();
var bodyParser = require('body-parser');

// Hyperledger Bridge
const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', 'network' ,'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

// Constants
const PORT = 8080;
const HOST = '0.0.0.0';

// use static file
app.use(express.static(path.join(__dirname, 'views')));

// configure app to use body-parser
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// main page routing
app.get('/', (req, res)=>{
    res.sendFile(__dirname + '/index.html');
})

async function cc_call(fn_name, args){
    
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);

    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('dolphins');

    var result;
    
    if(fn_name == 'addDiver')
    {
        i=args[0];
        n=args[1];
        bd=args[2];
        g=args[3];
        bt=args[4];
        result = await contract.submitTransaction('addDiver', i, n, bd, g, bt);
    }        
    else if( fn_name == 'addLevel')
    {
        i=args[0];
        ln=args[1];
        o=args[2];
        is=args[3];
        result = await contract.submitTransaction('addLevel', i, ln, o, is);
    }
    else if( fn_name == 'addCourse')
    {
        i=args[0];
        ln=args[1];
        c=args[2];
        result = await contract.submitTransaction('addCourse', i, ln, c);
    }
    else if( fn_name == 'addTestResult')
    {
        i=args[0];
        ln=args[1];
        s=args[2];
        result = await contract.submitTransaction('addTestResult', i, ln, s);
    }
    else if(fn_name == 'getLevel')
        result = await contract.evaluateTransaction('getLevel', args);
    else
        result = 'not supported function'

    return result;
}

// create diver
app.post('/diver', async(req, res)=>{
    const id = req.body.id;
    const name = req.body.name;  
    const bdate = req.body.bdate;
    const gender = req.body.gender;
    const btype = req.body.btype;
    console.log("add diver id: " + id);
    console.log("add diver name: " + name);
    console.log("add diver birthdate: " + bdate);
    console.log("add diver gender: " + gender);
    console.log("add diver blood type: " + btype);

    var args=[id, name, bdate, gender, btype];
    result = cc_call('addDiver', args)

    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})

// add level
app.post('/level', async(req, res)=>{
    const id = req.body.id;
    const levelname = req.body.levelname;
    const org = req.body.org;
    const instid = req.body.instid;
    console.log("add diver id: " + id);
    console.log("add level name: " + levelname);
    console.log("add organization: " + org);
    console.log("add instructor id: " + instid);

    var args=[id, levelname, org, instid];
    result = cc_call('addLevel', args)

    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})

// add course
app.post('/course', async(req, res)=>{
    const id = req.body.id;
    const levelname = req.body.levelname;
    const course = req.body.course;
    console.log("add diver id: " + id);
    console.log("add level name: " + levelname);
    console.log("add course: " + course);

    var args=[id, levelname, course];
    result = cc_call('addCourse', args)

    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})

// add test result
app.post('/test', async(req, res)=>{
    const id = req.body.id;
    const levelname = req.body.levelname;
    const status = req.body.status;
    console.log("add diver id: " + id);
    console.log("add level name: " + levelname);
    console.log("add status: " + status);

    var args=[id, levelname, status];
    result = cc_call('addTestResult', args)

    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})

// get level result
app.get('/diver', async (req,res)=>{
    const id = req.query.id;
    console.log("id: " + req.query.id);
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('dolphins');
    const result = await contract.evaluateTransaction('getLevel', id);
    const myobj = JSON.parse(result)
    res.status(200).json(myobj)
    // res.status(200).json(result)

});

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);